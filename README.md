# Random Word Generator
```
Random Word Generator
==========================
This module is a powerful tool for generating random words by scraping and wrapping the functionality of RandomWordGenerator.com. 

It allows easy access to the website's interface, enabling efficient random word generation for various applications.

Author: @z3ntl3
License: GNU
```
### Install module
``go get github.com/Z3NTL3/randomwordgenerator``

# Quickstart 
Look below for 2 examples

```go
package main

import (
	"github.com/Z3NTL3/randomwordgenerator"
	"fmt"
	"log"
)


func main() {
	client := randomwordgenerator.NewClient()
	err := client.SetQuantity(50); if err != nil {
		log.Fatal(err)
	}

	err = client.SetProxy("http://2.56.119.93:5074"); if err != nil {
		log.Fatal(err)
	}

	client.Initialize()

	words, err := client.GenerateWords(); if err != nil {
		log.Fatal(err)
	}
	for _, word := range words {
		fmt.Println(word)
	}
}

/*
OUTPUTL:
    extensiva
    desmajolariais
    dolarizaban
    funambulesco
    aglomeres
    chichinguastes
    envergabais
    jimar
    ensebare
    revalidaban
    gayad
    asentaban
    confiscara
    caricaturaramos
    arqueajes

    ...[30+ rows remaining]
*/
```

```go
package main

import (
	"github.com/Z3NTL3/randomwordgenerator"
	"context"
	"fmt"
	"log"
	"time"
)


func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	client := randomwordgenerator.WithContext(ctx)
	err := client.SetQuantity(50); if err != nil {
		log.Fatal(err)
	}

	client.Initialize()

	words, err := client.GenerateWords(); if err != nil {
		log.Fatal(err)
	}
	for _, word := range words {
		fmt.Println(word)
	}
}

/*
OUTPUT:
    contrastasen
    palmilla
    beborrotearemos
    recesariais
    saraveadas
    embarcarais
    desagraviara
    obsesionaran
    enjabonariamos
    curucutee

    ...[30+ rows remaining]
*/

```