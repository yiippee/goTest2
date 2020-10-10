// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2010 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// protoc-gen-micro is a plugin for the Google protocol buffer compiler to generate
// Go code.  Run it by building this program and putting it in your path with
// the name
// 	protoc-gen-micro
// That word 'micro' at the end becomes part of the option string set for the
// protocol compiler, so once the protocol compiler (protoc) is installed
// you can run
// 	protoc --micro_out=output_directory --go_out=output_directory input_directory/file.proto
// to generate go-micro code for the protocol defined by file.proto.
// With that input, the output will be written to
// 	output_directory/file.micro.go
//
// The generated code is documented in the package comment for
// the library.
//
// See the README and documentation for protocol buffers to learn more:
// 	https://developers.google.com/protocol-buffers/
package main

import (
	"io/ioutil"
	"os"

	"github.com/Allenxuxu/protoc-gen-mMicro/generator"
	_ "github.com/Allenxuxu/protoc-gen-mMicro/plugin/micro"
	"github.com/golang/protobuf/proto"
)

func main() {
	// Begin by allocating a generator. The request and response structures are stored there
	// so we can do error handling easily - the response structure contains the field to
	// report failure.
	g := generator.New()

	var err error

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
		panic(err)
	}
	////
	//fmt.Println("data:")
	// ioutil.WriteFile("./data.txt", data, os.ModePerm)

	// data, _ := ioutil.ReadFile("./data.txt")

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()

	g.GenerateAllFiles()

	// Send back the results.
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}

//
//var data = []byte(`10 10 116 101 115 116 46 112 114 111 116 111 26 8 8 3 16 11 24 4 34 0 122 148 9 10 10 116 101 115 116 46 112 114 111 116 111 18 4 116 10
//1 115 116 34 68 10 5 80 104 111 110 101 18 35 10 4 116 121 112 101 24 1 32 1 40 14 50 15 46 116 101 115 116 46 80 104 111 110 101 84 121 112 101 82 4 116 121 112 101 18 22 10 6 110 117 109 98 101 114 24 2 3
//2 1 40 9 82 6 110 117 109 98 101 114 34 101 10 6 80 101 114 115 111 110 18 14 10 2 105 100 24 1 32 1 40 5 82 2 105 100 18 18 10 4 110 97 109 101 24 2 32 1 40 9 82 4 110 97 109 101 18 35 10 6 112 104 111 110
// 101 115 24 3 32 3 40 11 50 11 46 116 101 115 116 46 80 104 111 110 101 82 6 112 104 111 110 101 115 18 18 10 4 116 101 115 116 24 4 32 1 40 5 82 4 116 101 115 116 34 53 10 11 67 111 110 116 97 99 116 66 11
//1 111 107 18 38 10 7 112 101 114 115 111 110 115 24 1 32 3 40 11 50 12 46 116 101 115 116 46 80 101 114 115 111 110 82 7 112 101 114 115 111 110 115 42 31 10 9 80 104 111 110 101 84 121 112 101 18 8 10 4 72
// 79 77 69 16 0 18 8 10 4 87 79 82 75 16 1 74 242 6 10 6 18 4 2 0 34 1 10 66 10 1 12 18 3 2 0 18 26 56 230 140 135 229 174 154 231 137 136 230 156 172 10 230 179 168 230 132 143 112 114 111 116 111 51 228 18
//4 142 112 114 111 116 111 50 231 154 132 229 134 153 230 179 149 230 156 137 228 186 155 228 184 141 229 144 140 10 10 52 10 1 2 18 3 5 0 13 26 42 229 140 133 229 144 141 239 188 140 233 128 154 232 191 135
// 112 114 111 116 111 99 231 148 159 230 136 144 230 151 182 103 111 230 150 135 228 187 182 230 151 182 10 10 63 10 2 5 0 18 4 9 0 12 1 26 51 230 137 139 230 156 186 231 177 187 229 158 139 10 230 158 154 2
//28 184 190 231 177 187 229 158 139 231 172 172 228 184 128 228 184 170 229 173 151 230 174 181 229 191 133 233 161 187 228 184 186 48 10 10 10 10 3 5 0 1 18 3 9 5 14 10 11 10 4 5 0 2 0 18 3 10 4 13 10 12 10
// 5 5 0 2 0 1 18 3 10 4 8 10 12 10 5 5 0 2 0 2 18 3 10 11 12 10 11 10 4 5 0 2 1 18 3 11 4 13 10 12 10 5 5 0 2 1 1 18 3 11 4 8 10 12 10 5 5 0 2 1 2 18 3 11 11 12 10 19 10 2 4 0 18 4 15 0 18 1 26 7 230 137 139
// 230 156 186 10 10 10 10 3 4 0 1 18 3 15 8 13 10 11 10 4 4 0 2 0 18 3 16 4 23 10 12 10 5 4 0 2 0 6 18 3 16 4 13 10 12 10 5 4 0 2 0 1 18 3 16 14 18 10 12 10 5 4 0 2 0 3 18 3 16 21 22 10 11 10 4 4 0 2 1 18 3
//17 4 22 10 12 10 5 4 0 2 1 5 18 3 17 4 10 10 12 10 5 4 0 2 1 1 18 3 17 11 17 10 12 10 5 4 0 2 1 3 18 3 17 20 21 10 16 10 2 4 1 18 4 21 0 29 1 26 4 228 186 186 10 10 10 10 3 4 1 1 18 3 21 8 14 10 44 10 4 4 1
// 2 0 18 3 23 4 17 26 31 229 144 142 233 157 162 231 154 132 230 149 176 229 173 151 232 161 168 231 164 186 230 160 135 232 175 134 229 143 183 10 10 12 10 5 4 1 2 0 5 18 3 23 4 9 10 12 10 5 4 1 2 0 1 18 3
//23 10 12 10 12 10 5 4 1 2 0 3 18 3 23 15 16 10 11 10 4 4 1 2 1 18 3 24 4 20 10 12 10 5 4 1 2 1 5 18 3 24 4 10 10 12 10 5 4 1 2 1 1 18 3 24 11 15 10 12 10 5 4 1 2 1 3 18 3 24 18 19 10 59 10 4 4 1 2 2 18 3 27
// 4 30 26 46 114 101 112 101 97 116 101 100 232 161 168 231 164 186 229 143 175 233 135 141 229 164 141 10 229 143 175 228 187 165 230 156 137 229 164 154 228 184 170 230 137 139 230 156 186 10 10 12 10 5 4
//1 2 2 4 18 3 27 4 12 10 12 10 5 4 1 2 2 6 18 3 27 13 18 10 12 10 5 4 1 2 2 1 18 3 27 19 25 10 12 10 5 4 1 2 2 3 18 3 27 28 29 10 11 10 4 4 1 2 3 18 3 28 4 19 10 12 10 5 4 1 2 3 5 18 3 28 4 9 10 12 10 5 4 1
//2 3 1 18 3 28 10 14 10 12 10 5 4 1 2 3 3 18 3 28 17 18 10 22 10 2 4 2 18 4 32 0 34 1 26 10 232 129 148 231 179 187 231 176 191 10 10 10 10 3 4 2 1 18 3 32 8 19 10 11 10 4 4 2 2 0 18 3 33 4 32 10 12 10 5 4 2
// 2 0 4 18 3 33 4 12 10 12 10 5 4 2 2 0 6 18 3 33 13 19 10 12 10 5 4 2 2 0 1 18 3 33 20 27 10 12 10 5 4 2 2 0 3 18 3 33 30 31 98 6 112 114 111 116 111 51`)
