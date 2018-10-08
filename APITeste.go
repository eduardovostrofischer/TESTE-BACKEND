package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gedex/bp3d"
	"github.com/gorilla/mux"
)

type Product struct {
	ID          string  `json:"id"`
	QTD         int     `json:"qtd"`
	Largura     float64 `json:"largura"`
	Altura      float64 `json:"altura"`
	Comprimento float64 `json:"comprimento"`
	Peso        float64 `json:"peso"`
}
type RetBox struct {
	Largura     float64 `json:"largura"`
	Altura      float64 `json:"altura"`
	Comprimento float64 `json:"comprimento"`
	Peso        float64 `json:"peso"`

	RetItems []RetItem `json: "it"`
}
type RetItem struct {
	ID  string `json:"id"`
	QTD int    `json:"qtd"`
}

//Init products var as a slice product struct
var products []Product

//Get all products
func getproducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

//Get single product
func getproduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//Loop throught products and +find with in
	for _, item := range products {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Product{})
}

// Remove product
func removeproduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range products {
		if item.ID == params["id"] {
			products[index].QTD--
			if products[index].QTD <= 0 {
				products = append(products[:index], products[index+1:]...)
				break
			}
		}
	}
	json.NewEncoder(w).Encode(products)
}

// Insert product, se produto ja existe apenas aumenta a quantidade
func insertproduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productexists := false
	var product Product
	_ = json.NewDecoder(r.Body).Decode(&product)
	for index, item := range products {
		if item.ID == product.ID {
			products[index].QTD += product.QTD
			productexists = true
			break
		}
	}
	// ID do novo produto é gerado por um número aleatório. Em produção obviamente não seria uma escolha
	// adequada, uma conexão com o banco de dados é necessária para associar um id correto.
	if !productexists {
		product.ID = strconv.Itoa(rand.Intn(10000000)) //Mock -not safe
		products = append(products, product)
		json.NewEncoder(w).Encode(product)
	} else {
		json.NewEncoder(w).Encode(products)
	}
}

// Devolve as caixas com os itens dentro delas
func getboxes(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	p := bp3d.NewPacker()
	params := mux.Vars(r)
	var transpLargura float64
	var transpAltura float64
	var transpComprimento float64
	var transpPeso float64

	// Em um cenário real recuperaríamos valores de altura,largura,comprimento e peso da transportadora
	// de um banco de dados, para esse teste usei um switch case.
	switch params["transp"] {
	case "correios":
		transpAltura = 105
		transpLargura = 105
		transpComprimento = 105
		transpPeso = 30
	case "jadlog":
		transpAltura = 100
		transpLargura = 105
		transpComprimento = 181
		transpPeso = 150
	case "viabrasil":
		transpAltura = 200
		transpLargura = 200
		transpComprimento = 200
		transpPeso = 200
	}

	for r, prod := range products {
		// Verifica se dimensões e peso de um item são adequados
		if prod.Largura > transpLargura || prod.Altura > transpAltura || prod.Comprimento > transpComprimento || prod.Peso > transpPeso {
			json.NewEncoder(w).Encode("item de id=" + prod.ID + " ultrapassa os limites da transportadora Correios")
			return
		}
		// Adicionei uma caixa para cada item, é menos custoso adicionar mais caixas do que tentar estimar um
		// número de caixas necessárias, removeremos as caixas vazias depois
		for i := 0; i < prod.QTD; i++ {
			p.AddBin(bp3d.NewBin(params["transp"]+strconv.Itoa(r)+" "+strconv.Itoa(i), transpLargura, transpAltura, transpComprimento, transpPeso))
			p.AddItem(bp3d.NewItem(prod.ID, prod.Largura, prod.Altura, prod.Comprimento, prod.Peso))
		}
	}

	// Realiza o 3d binpacking
	if err := p.Pack(); err != nil {
		log.Fatal(err)
	}

	// Soma do peso total da caixa
	var totalWeight float64
	for index, bins := range p.Bins {
		totalWeight = 0
		for _, ib := range bins.Items {
			totalWeight += ib.Weight
		}
		p.Bins[index].Weight = totalWeight
	}

	// Calcula as Dimensões da caixa,
	for index, boxes := range p.Bins {
		var maxw float64
		var maxh float64
		var maxd float64
		for _, item := range boxes.Items {
			str := item.RotationType.String()
			str = str[13:16]
			// TODO: Esse switch case pode ser encapsulado em uma função.
			switch str {
			case "WHD":
				if maxw < item.Width+item.Position[0] { //W
					maxw = item.Width + item.Position[0]
				}
				if maxh < item.Height+item.Position[1] { //H
					maxh = item.Height + item.Position[1]
				}
				if maxd < item.Depth+item.Position[2] { //D
					maxd = item.Depth + item.Position[2]
				}
			case "HWD":
				if maxw < item.Height+item.Position[1] { //H
					maxw = item.Height + item.Position[1]
				}
				if maxh < item.Width+item.Position[0] { //W
					maxh = item.Width + item.Position[0]
				}
				if maxd < item.Depth+item.Position[2] { //D
					maxd = item.Depth + item.Position[2]
				}
			case "HDW":
				if maxw < item.Height+item.Position[1] { //H
					maxw = item.Height + item.Position[1]
				}
				if maxh < item.Depth+item.Position[2] { //D
					maxh = item.Depth + item.Position[2]
				}
				if maxd < item.Width+item.Position[0] { //W
					maxd = item.Width + item.Position[0]
				}
			case "DHW":
				if maxw < item.Depth+item.Position[2] { //D
					maxw = item.Depth + item.Position[2]
				}
				if maxh < item.Height+item.Position[1] { //H
					maxh = item.Height + item.Position[1]
				}
				if maxd < item.Width+item.Position[0] { //W
					maxd = item.Width + item.Position[0]
				}
			case "DWH":
				if maxw < item.Depth+item.Position[2] { //D
					maxw = item.Depth + item.Position[2]
				}
				if maxh < item.Width+item.Position[0] { //W
					maxh = item.Width + item.Position[0]
				}
				if maxd < item.Height+item.Position[1] { //H
					maxd = item.Height + item.Position[1]
				}
			case "WDH":
				if maxw < item.Width+item.Position[0] { //W
					maxw = item.Width + item.Position[0]
				}
				if maxh < item.Depth+item.Position[2] { //D
					maxh = item.Depth + item.Position[2]
				}
				if maxd < item.Height+item.Position[1] { //H
					maxd = item.Height + item.Position[1]
				}
			}
		}
		p.Bins[index].Width = maxw
		p.Bins[index].Height = maxh
		p.Bins[index].Depth = maxd
	}

	// Remove as caixas vazias
	for index, b := range p.Bins {
		if b.Weight == 0 {
			p.Bins = p.Bins[0:index]
			break
		}
	}
	// Transformações para deixar a resposta no formato correto
	var ret RetBox
	var reti RetItem
	var retBoxes []RetBox
	for index, b := range p.Bins {
		ret.Largura = b.Width
		ret.Altura = b.Height
		ret.Comprimento = b.Depth
		ret.Peso = b.Weight
		retBoxes = append(retBoxes, ret)
		for _, i := range b.Items {
			reti.ID = i.Name
			reti.QTD = 1
			retBoxes[index].RetItems = append(retBoxes[index].RetItems, reti)
		}
	}

	// Une e quantifica as caixas de mesmo id
	for index, b := range p.Bins {
		m := make(map[string]int)
		for _, i := range b.Items {
			m[i.Name]++
		}
		retBoxes[index].RetItems = retBoxes[index].RetItems[:0]
		for str := range m {
			reti.ID = str
			reti.QTD = m[str]
			retBoxes[index].RetItems = append(retBoxes[index].RetItems, reti)
		}
	}

	displayPacked(p.Bins)
	json.NewEncoder(w).Encode(retBoxes)
}

// Print dos itens empacotados
func displayPacked(bins []*bp3d.Bin) {
	for _, b := range bins {
		fmt.Println(b)
		fmt.Println(" packed items:")
		for _, i := range b.Items {
			fmt.Println("  ", i)
		}
	}
}

func main() {
	//Init Router
	r := mux.NewRouter()

	products = append(products, Product{ID: "1", QTD: 1, Altura: 10, Largura: 10, Comprimento: 10, Peso: 29})
	products = append(products, Product{ID: "5", QTD: 3, Altura: 15, Largura: 15, Comprimento: 20, Peso: 5})

	//Route Handlers / Endpoints
	r.HandleFunc("/api/products", getproducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", getproduct).Methods("GET")
	r.HandleFunc("/api/products", insertproduct).Methods("POST")
	r.HandleFunc("/api/boxes/{transp}", getboxes).Methods("GET")
	r.HandleFunc("/api/products/{id}", removeproduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
