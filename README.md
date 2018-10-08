# TESTE-BACKEND

Aqui o github do teste de back-end da empresa melhor envio.

Todos os requisitos foram atendidos.

Dependências:
Instale Go em seu computador

sudo apt-get install golang-go.

Instale a biblioteca github.com/gorilla/mux 

go get -u github.com/gorilla/mux.

Descrição do problema:
O problema apresentado foi 3d bin packing com pesos. É um problema NP-Difícil, logo soluções para ele envolvem complexidades de tempo maiores do que polinomiais. Assim ele se torna muito custoso mesmo para entradas pequenas. 

Instruções:
Depois de ter instalado go e a biblioteca gorila mux.

Adicione gedex/bp3d dentro de sua pasta github.com dentro de src em sua instalação go.

Coloque APITeste.go dentro de uma pasta APITeste dentro de src em sua instalação go.

Compile e execute ApiTeste.go

Recomendo teste com POSTMAN.

Referências:
https://github.com/gedex/bp3d. 
Código de 3d-binpacking foi baseado neste porém com pequenas alterações como correções de pesos e variáveis adicionais para as caixas.

https://github.com/bom-d-van/binpacking/blob/master/erick_dube_507-034.pdf.
O algoritmo solução é baseado nesse artigo, é um algoritmo de best fit, onde se procura o menor desperdício de volume possível.

5 endpoints na api

getproducts     /api/products               GET
  Mostra todos produtos do carrinho
// O programa ja inicia com 4 itens no carrinho

getproduct      /api/products/{id}          GET
  Mostra um produto específico do carrinho

insertproduct   /api/products               POST
  adiciona um produto no carrinho

removeproduct   /api/products/{id}          DELETE
  remove um produto do carrinho

getboxes        /api/boxes/{transp}         GET
  calcula as caixas necessárias, dimensões e peso e os itens que a caixa contém
  {transp} = "correios","jadlog","viabrasil"
  
