package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"net/url"
	"nvnoskov/presentation-builder/models"
	template "nvnoskov/presentation-builder/templates"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
	"github.com/savsgio/atreugo/v9"
	"github.com/savsgio/go-logger"
	"github.com/thedevsaddam/govalidator"
	"github.com/valyala/fasthttp"
)

func init() {
	logger.SetLevel(logger.ERROR)
}

func main() {

	// Database initialization
	models.Connect()

	config := &atreugo.Config{
		Addr: "0.0.0.0:8000",
	}
	server := atreugo.New(config)

	server.Static("/static", "static")
	// Register authentication middleware at first of all
	server.UseBefore(authMiddleware)

	// Register index route
	server.Path("GET", "/", func(ctx *atreugo.RequestCtx) error {
		return ctx.RedirectResponse("/presentations", ctx.Response.StatusCode())
	})

	// Register login route
	server.Path("GET", "/login", func(ctx *atreugo.RequestCtx) error {
		buffer := new(bytes.Buffer)
		template.Login(buffer)

		return ctx.HTTPResponseBytes(buffer.Bytes())
	})

	// Register presentations route
	server.Path("GET", "/presentations", func(ctx *atreugo.RequestCtx) error {
		user := ctx.UserValue("user").(*userCredential)

		presentations := []models.Presentation{}

		err := models.DB.Select(&presentations, "SELECT * FROM presentations WHERE slug=$1 and draft=0", user.Slug)
		if err != nil {
			logger.Error(err)
		}
		fmt.Printf("pres: %+v %s", presentations, user.Slug)

		buffer := new(bytes.Buffer)
		template.Presentations(user.Slug, presentations, buffer)

		return ctx.HTTPResponseBytes(buffer.Bytes())
	})

	// Register login route
	server.Path("POST", "/login", func(ctx *atreugo.RequestCtx) error {
		qMail := ctx.PostArgs().Peek("email")

		jwtCookie := ctx.Request.Header.Cookie("atreugo_jwt")

		if len(jwtCookie) == 0 {
			token, tokenString, expireAt := generateToken(qMail)

			// Set cookie for domain
			cookie := fasthttp.AcquireCookie()
			defer fasthttp.ReleaseCookie(cookie)

			cookie.SetKey("atreugo_jwt")
			cookie.SetValue(tokenString)
			cookie.SetExpire(expireAt)
			ctx.Response.Header.SetCookie(cookie)
			_, err := models.DB.Exec(`INSERT INTO users(email,slug) 
			SELECT $1, $2 
			WHERE NOT EXISTS(SELECT 1 FROM users WHERE email = $1 AND slug = $2);`, token.Email, token.Slug)
			if err != nil {
				logger.Error(err)
			}

		}
		return ctx.RedirectResponse("/", ctx.Response.StatusCode())
	})
	server.Path("POST", "/upload", func(ctx *atreugo.RequestCtx) error {
		user := ctx.UserValue("user").(*userCredential)

		dirPath := fmt.Sprintf("static/presentation/%s/", user.Slug)
		filePath := dirPath + "file.pdf"

		fh, err := ctx.FormFile("file")
		if err != nil {
			logger.Error(err)
		}

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			os.Mkdir(dirPath, 0755)
		}
		if err := fasthttp.SaveMultipartFile(fh, filePath); err != nil {
			logger.Error(err)
		}
		f, err := os.Open(filePath)
		if err != nil {
			logger.Fatal(err)
		}
		defer f.Close()

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			logger.Fatal(err)
		}

		hash := fmt.Sprintf("%x", h.Sum(nil))

		shortHash := hash[:8]
		var count int
		err = models.DB.Get(&count, "SELECT COUNT(file) FROM presentations WHERE file=$1 AND slug = $2", shortHash, user.Slug)
		if err != nil {
			logger.Error(err)
		}
		logger.Info("COUNT:", count, shortHash)

		if count > 0 {
			return ctx.RedirectResponse("/edit/"+shortHash, ctx.Response.StatusCode())
		}

		doc, err := fitz.New(filePath)
		if err != nil {
			logger.Error(err)
		}
		defer doc.Close()
		if _, err := os.Stat(dirPath + shortHash); os.IsNotExist(err) {
			os.Mkdir(dirPath+shortHash, 0755)
		}
		// Extract pages as images
		for n := 0; n < doc.NumPage(); n++ {
			img, err := doc.Image(n)
			if err != nil {
				logger.Error(err)
			}

			f, err := os.Create(filepath.Join(dirPath, shortHash, fmt.Sprintf("page-%02d.jpg", n)))
			if err != nil {
				logger.Error(err)
			}

			err = jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
			if err != nil {
				logger.Error(err)
			}

			f.Close()
		}

		_, err = models.DB.Exec(`INSERT INTO presentations(file,slug,draft,pages, author, name, description) 
			SELECT $1, $2, 0, $3, "", "" , ""
			WHERE NOT EXISTS(SELECT 1 FROM presentations WHERE file = $1 AND slug = $2);`, hash[:8], user.Slug, doc.NumPage())
		if err != nil {
			logger.Error(err)
		}
		os.Remove(filePath)

		return ctx.RedirectResponse("/edit/"+shortHash, ctx.Response.StatusCode())
	})
	server.Path("GET", "/upload", func(ctx *atreugo.RequestCtx) error {

		user := ctx.UserValue("user").(*userCredential)
		buffer := new(bytes.Buffer)

		// check active pdf document
		if _, err := os.Stat(fmt.Sprintf("static/presentation/%s/file.pdf", user.Slug)); os.IsNotExist(err) {

			// If not exists we show upload form
			template.Upload(buffer)
			return ctx.HTTPResponseBytes(buffer.Bytes())
		} else {
			// redirect to editing
			return ctx.RedirectResponse("/edit", ctx.Response.StatusCode())
		}

	})

	server.Path("GET", "/preview/:file?", previewPresentation)
	server.Path("GET", "/edit/:file?", editPresentation)
	server.Path("POST", "/edit/:file?", editPresentation)
	server.Path("POST", "/audio", saveAudio)

	// Run
	err := server.ListenAndServe()
	if err != nil {
		logger.Error(err)
	}
}

func saveAudio(ctx *atreugo.RequestCtx) error {
	user := ctx.UserValue("user").(*userCredential)

	var dat map[string]interface{}
	if err := json.Unmarshal(ctx.Request.Body(), &dat); err != nil {
		logger.Error(err)
	}
	if user.Slug != dat["slug"].(string) {
		return ctx.HTTPResponse("access denied", 403)
	}

	dec, err := base64.StdEncoding.DecodeString(dat["message"].(string))
	if err != nil {
		logger.Error(err)
	}
	f, err := os.Create(fmt.Sprintf("static/presentation/%s/%s/page-%02.0f.wav", user.Slug, dat["file"].(string), dat["page"].(float64)))
	if err != nil {
		logger.Error(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		logger.Error(err)
	}
	if err := f.Sync(); err != nil {
		logger.Error(err)
	}
	return ctx.HTTPResponse("ok", 201)

}
func previewPresentation(ctx *atreugo.RequestCtx) error {

	user := ctx.UserValue("user").(*userCredential)
	buffer := new(bytes.Buffer)
	presentation := models.Presentation{}
	var err error
	file := ctx.UserValue("file") // Could be nil

	err = models.DB.Get(&presentation, "SELECT * FROM presentations WHERE file=$1 and slug=$2", file, user.Slug)

	if err != nil {
		logger.Error(err)
	}
	var entry models.Entry
	var entries []models.Entry
	for i := 0; i < presentation.Pages; i++ {
		entry = models.Entry{
			Image: fmt.Sprintf("/static/presentation/%s/%s/page-%02d.jpg", user.Slug, presentation.File, i),
		}
		_, err := os.Stat(fmt.Sprintf("static/presentation/%s/%s/page-%02d.wav", user.Slug, presentation.File, i))

		if err == nil {
			entry.Audio = fmt.Sprintf("/static/presentation/%s/%s/page-%02d.wav", user.Slug, presentation.File, i)
		}
		logger.Error(err)
		logger.Info(entry)
		entries = append(entries, entry)
	}
	presentation.Entries = entries
	template.Preview(user.Slug, presentation, buffer)
	return ctx.HTTPResponseBytes(buffer.Bytes())
}
func editPresentation(ctx *atreugo.RequestCtx) error {

	user := ctx.UserValue("user").(*userCredential)
	buffer := new(bytes.Buffer)
	presentation := models.Presentation{}
	var err error
	var e url.Values
	file := ctx.UserValue("file") // Could be nil

	err = models.DB.Get(&presentation, "SELECT * FROM presentations WHERE file=$1 and slug=$2", file, user.Slug)

	if ctx.IsPost() {
		presentation.Name = string(ctx.PostArgs().Peek("name"))
		presentation.Author = string(ctx.PostArgs().Peek("author"))
		presentation.Description = string(ctx.PostArgs().Peek("description"))

		rules := govalidator.MapData{
			"author":      []string{"required", "max:255"},
			"name":        []string{"required", "min:4", "max:255"},
			"description": []string{"required"},
		}

		messages := govalidator.MapData{
			"name": []string{"required:Enter the name of your presentation"},
			// "author":      []string{"required:Enter the author of your presentation"},
			"description": []string{"required:Enter the description of your presentation"},
		}
		opts := govalidator.Options{
			Data:            &presentation,
			Rules:           rules,    // rules map
			Messages:        messages, // custom message map (Optional)
			RequiredDefault: true,     // all the field to be pass the rules
		}
		v := govalidator.New(opts)
		e = v.ValidateStruct()

		if len(e) == 0 {
			// Save presentation record

			_, err = models.DB.NamedExec(`UPDATE presentations 
			SET name=:name, author=:author, description=:description			
			WHERE file=:file and slug=:slug`, map[string]interface{}{
				"name":        presentation.Name,
				"author":      presentation.Author,
				"description": presentation.Description,
				"slug":        user.Slug,
				"file":        file,
			})

			if err != nil {
				logger.Error(err)
			}
			return ctx.RedirectResponse("/presentations", ctx.Response.StatusCode())
		}
	}
	template.Edit(e, presentation, buffer)
	return ctx.HTTPResponseBytes(buffer.Bytes())

}
