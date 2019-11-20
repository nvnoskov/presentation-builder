// Code generated by hero.
// source: /home/strannik/www/presentation-builder/templates/preview.html
// DO NOT EDIT!
package template

import (
	"bytes"
	"nvnoskov/presentation-builder/models"

	"github.com/shiyanhui/hero"
)

func Preview(slug string, presentation models.Presentation, buffer *bytes.Buffer) {
	buffer.WriteString(`<html>

<head>

    <!-- Bootstrap core CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
        integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
    <link rel="stylesheet" href="/static/style.css">

    <link rel="stylesheet" href="/static/reveal/css/reveal.css">
    <link rel="stylesheet" href="/static/reveal/css/theme/white.css">
</head>

<body>
    `)
	buffer.WriteString(`
    <script>
        window.currentSlug = '`)
	hero.EscapeHTML(slug, buffer)
	buffer.WriteString(`';    
        window.currentUrl = '`)
	hero.EscapeHTML(presentation.File, buffer)
	buffer.WriteString(`';    
    </script>
    <style>
        .slides audio {
            position: fixed;
            top:0;
            right: 0;
        }
    </style>
    <nav class="navbar fixed-bottom navbar-light bg-light">
            <button class="btn btn-danger" id="record">Record</button>
            <button class="btn btn-info" id="stop" disabled>Stop</button>
            <button class="btn btn-info" id="play" disabled>Play</button>
            <button class="btn btn-primary" id="save" disabled>Save</button>
    </nav>
    <div class="reveal">
        <div class="slides">
            `)
	for _, entry := range presentation.Entries {
		buffer.WriteString(`
                <section data-background-image="`)
		hero.EscapeHTML(entry.Image, buffer)
		buffer.WriteString(`">
                    `)
		if entry.Audio != "" {
			buffer.WriteString(`
                        <audio controls>
                            <source src="`)
			hero.EscapeHTML(entry.Audio, buffer)
			buffer.WriteString(`" type="audio/wav">                        
                        Your browser does not support the audio element.
                        </audio> 
                        `)
		}
		buffer.WriteString(`
                </section>
            `)
	}
	buffer.WriteString(`            
            
        </div>
    </div>
    
    <script src="/static/reveal/js/reveal.js"></script>
    <script>
        Reveal.initialize({
            autoPlayMedia: true,
            dependencies: [
                {
                    src: '/static/reveal/js/audio.js', async: true 
                },
            ]
        });
    </script>
`)

	buffer.WriteString(`

</body>

</html>`)

}
