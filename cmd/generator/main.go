// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"os"

	"github.com/googleapis/generator/internal/generator"
)

func main() {
	ctx := context.Background()
	log.Println("Invoking generator with arguements:", strings.Join(os.Args[1:], " "))
	if err := generator.Run(ctx, os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}
