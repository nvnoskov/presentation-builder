# CHALLENGE 3 - PRESENTATION BUILDER

## Basic Specs 
1. **DONE** Should accept a PDF document and split it into an image per page. 
2. **TODO** The PDF document should have at least one page and no more than twenty. 
3. **DONE** Upon having the PDF split into an image per page, the user should be able to record audio for each page. Audio is optional per page.
4. The presentation itself should have a name (not to exceed 255 characters) and a description (unlimited size). The description should allow HTML. 
5. **DONE** Each page within the presentation should consist of an image (required) and an audio file (optional). Each page within the presentation should also store its sort order within the presentation. Each page should be sorted by default in the order it appears in the original PDF document. Each page within the presentation should also know the length of the audio (in seconds) if the audio is present. 
6. **DONE** It should be possible to edit a presentation. The end user should be able to edit the name and the HTML description.  
7. **PARTIAL** It should be possible to edit a presentation page. The end user should be able to change the image and to replace and/or delete the audio  associated with the page. The end user should also be able to re-order a presentation page within the presentation. Additionally, the end user should be able to insert a new presentation page within the presentation by providing a new image. They will also be able to record audio for that page. The addition of pages should only be allowed if the presentation has less than 20 pages. 
8. **TODO** The end user should be able to export their presentation. An exported presentation should be in the form of a zip file. The zip file should contain all images and audio files that make up the presentation. Additionally, a manifest file should be present that describes how the images and audio files are related within the presentation. This is to support importing the presentation.
9. **TODO** The manifest file should be structured as a JSON file and should include the following information: 
    1. Author 
    2. Date 
    3. Presentation Name 
    4. Presentation Description 
    5. An entry per page, where each entry contains: 
        1. The sort order of the page
        2. The name of the audio file within the zip if the audio file is present 
        3. The name of the image file within the zip 
        4. The length of the audio file in seconds if the audio file is present


### Tech stack

- Web application based on  [Atreugo](https://github.com/savsgio/atreugo/)
- Template engine [Hero](https://github.com/shiyanhui/hero)
- Presentation library [Reveal.js](https://github.com/hakimel/reveal.js)

## Run
```
git clone  https://github.com/nvnoskov/presentation-builder.git
cd presentation-builder
go mod download
go run *.go
```
## Build

#### Hero template compilation
```
hero -source=./templates
```
#### Docker image
```
docker build -t presentation-builder .
```
#### Run
```
docker run -it --rm -p8001:8000 presentation-builder
```

### Improvments 
- KeyBindings!


### TODO

 - More abstractions (file store, databases, etc.)
 - validations
 - Patterns (Repository)
 - MP3/OGG-encoding



### Benchmark 

Test template renderer performance on for the login page
```
wrk -t12 -c400 -d30s http://localhost:8000/login
Running 30s test @ http://localhost:8000/login
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.77ms    7.01ms  86.97ms   87.07%
    Req/Sec     8.41k     1.49k   23.33k    69.86%
  3014556 requests in 30.10s, 8.45GB read
Requests/sec: 100153.09
Transfer/sec:    287.59MB

```
Machine:
- CPU: AMD® Fx-8320e eight-core processor × 8 
- Memory: 15.5 GiB