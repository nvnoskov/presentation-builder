// Code generated by hero.
// source: /home/strannik/www/presentation-builder/templates/presenations.html
// DO NOT EDIT!
package template

import (
	"bytes"
	"nvnoskov/presentation-builder/models"

	"github.com/shiyanhui/hero"
)

func Presentations(slug string, presentations []models.Presentation, buffer *bytes.Buffer) {
	buffer.WriteString(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="">
    <title>Presentation builder</title>

    <!-- Bootstrap core CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
    <link rel="stylesheet" href="/static/style.css">
   
  </head>
  <body>
        <nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
                <a class="navbar-brand" href="/">Presentation builder</a>
                <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
                  <span class="navbar-toggler-icon"></span>
                </button>
                <div class="collapse navbar-collapse" id="navbarCollapse">
                  <ul class="navbar-nav mr-auto">
                    <li class="nav-item active">
                      <a class="nav-link" href="/">Home</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="/presentations">My presentations</a>
                    </li>                    
                    <li class="nav-item">
                      <a class="nav-link" href="/upload">New upload</a>
                    </li>                    
                  </ul>
                  
                </div>
              </nav>
              
              <main role="main" class="container">
        `)
	buffer.WriteString(`
    <style>
        .col{
            margin-bottom:10px;
        }
    </style>
    <div class="jumbotron">
            <h1>Presentations</h1>
            <div class="row align-items-center">
            `)
	for _, presentation := range presentations {
		buffer.WriteString(`
                <div class="col">
                    `)
		buffer.WriteString(`<div class="card" style="width: 18rem;">
    <img src="/static/presentation/`)
		hero.EscapeHTML(slug, buffer)
		buffer.WriteString(`/`)
		hero.EscapeHTML(presentation.File, buffer)
		buffer.WriteString(`/page-00.jpg" class="card-img-top" alt="`)
		hero.EscapeHTML(presentation.Name, buffer)
		buffer.WriteString(`">
    <div class="card-body">
        <h5 class="card-title">`)
		hero.EscapeHTML(presentation.Name, buffer)
		buffer.WriteString(`</h5>
        <p class="card-text">`)
		hero.EscapeHTML(presentation.Description, buffer)
		buffer.WriteString(`</p>
        <a href="/edit/`)
		hero.EscapeHTML(presentation.File, buffer)
		buffer.WriteString(`" class="btn btn-primary">Edit</a>
        <a href="/preview/`)
		hero.EscapeHTML(presentation.File, buffer)
		buffer.WriteString(`" class="btn btn-success">Preview</a>
        <a href="/delete/`)
		hero.EscapeHTML(presentation.File, buffer)
		buffer.WriteString(`" class="btn btn-light">Remove</a>
    </div>
    </div>`)
		buffer.WriteString(`
                </div>
            `)
	}
	buffer.WriteString(`
            </div>
            <br>
            <a href="/upload" class="btn btn-block btn-primary btn-lg">Upload</a>
          </div>

    
`)

	buffer.WriteString(`
        </main>
    <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
</body>
</html>
`)

}
