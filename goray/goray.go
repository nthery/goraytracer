/*
Copyright (c) 2013 Nicolas Thery <nthery@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
 Minimal ray tracer command line interface.

 Parse a JSON-encoded scene, render it, and write it to a PNG file.

 TODO: Replace JSON with better input format
*/
package main

import (
	"encoding/json"
	"flag"
	"github.com/nthery/goraytracer/raytracer"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
)

var (
	njobs      = flag.Int("j", 1, "# of parallel jobs")
	infile     = flag.String("i", "", "input file")
	outfile    = flag.String("o", "", "output file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	loop       = flag.Int("l", 1, "# of rendering loop (for profiling)")
)

func main() {
	flag.Parse()
	if *infile == "" {
		log.Fatalf("no input file")
	}
	if *outfile == "" {
		log.Fatalf("no output file")
	}
	if *njobs < 1 {
		*njobs = 1
	}

	fin, err := os.Open(*infile)
	if err != nil {
		log.Fatalf("can not open input file: %v\n", err)
	}
	defer fin.Close()

	fout, err := os.Create(*outfile)
	if err != nil {
		log.Fatalf("can not create output file: %v\n", err)
	}
	defer fout.Close()

	in, err := ioutil.ReadAll(fin)
	if err != nil {
		log.Fatalf("can not read input file: %v\n", err)
	}

	var scene raytracer.Scene
	err = json.Unmarshal(in, &scene)
	if err != nil {
		log.Fatalf("can not parse input file: %v\n", err)
	}

	img, err := renderScene(&scene)
	if err != nil {
		log.Fatalf("can not render scene: %v\n", err)
	}

	err = png.Encode(fout, img)
	if err != nil {
		log.Fatalf("can not generate PNG output: %v\n", err)
	}
}

func renderScene(s *raytracer.Scene) (*image.RGBA, error) {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("can not create profile file: ", err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	for i := 0; i < *loop-1; i++ {
		s.Render(*njobs)
	}
	return s.Render(*njobs)
}
