package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

// Struct utilizada para exibir dados no template
type Names struct {
	Id    int
	Name  string
	Email string
}

// Função dbConn, abrea a conexão com o banco de dados
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "estevamnet"
	dbPass := "archlinux"
	dbName := "crud"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// A variável tmpl renderiza todos os templates da pasta 'tmpl' independente da extensão
var tmpl = template.Must(template.ParseGlob("tmpl/*"))

// Função  usada para renderizar o arquivo Index
func Index(w http.ResponseWriter, r *http.Request) {
	// Abre a conexão com obanco de dados utilizando a função dbConn()
	db := dbConn()
	// Reqliza a consulta com obanco de dados e trata erros
	selDB, err := db.Query("SELECT * FROM names ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	// Monta a struct
	n := Names{}
	// Monta um Array para guardar os valores da struct
	res := []Names{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		// Armazena os valores em variáveis
		var id int
		var name, email string

		// Faz o scan do select
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email

		// Junta a struct com array
		res = append(res, n)
	}

	// abre a página Index e exibe todos os regisreados na tela
	tmpl.ExecuteTemplate(w, "Index", res)

	// Frcha a conexão
	defer db.Close()
}

func main() {
	// Exibe a mensagem que os servidor foi iniciado
	log.Println("Server started on http://localhost:9000")

	// Geremcia as URLs
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)

	// Ações
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	// Inicia o servidor na porta 9000
	http.ListenAndServe(":9000", nil)
}

// Funçção Show
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	// Pega o Id do parâmetro da URL
	nId := r.URL.Query().Get("id")
	// Usa o Id para fazer a consulta e tratar erros
	selDB, err := db.Query("SELECT  * FROM  names WHERE  id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	// Monta a struct paras er utilizada no template
	n := Names{}
	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		// armazena os valores em variáveis
		var id int
		var name, email string

		// Faz o scan do SELECT
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		// Envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email
	}

	// Mostra o template
	tmpl.ExecuteTemplate(w, "Show", n)
	// fecha a conexão
	defer db.Close()
}

// Função new
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Função Edit
func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	//Pega o parâmetro do Id da URL
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM names WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	n := Names{}

	for selDB.Next() {
		var id int
		var name, email string

		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		n.Id = id
		n.Name = name
		n.Email = email
	}

	tmpl.ExecuteTemplate(w, "Edit", n)

	defer db.Close()
}

// Fumção Insert, insere valores no banco de dados
func Insert(w http.ResponseWriter, r *http.Request) {
	// Abre a conexão
	db := dbConn()
	// Verifica o METHOD do formulário passado
	if r.Method == "POST" {
		// Pega os campos do formulário
		name := r.FormValue("name")
		email := r.FormValue("email")
		// Prepara o SQL e verifica erros
		insForm, err := db.Prepare("INSERT INTO names(name, email) VALUES (?,?)")
		if err != nil {
			panic(err.Error())
		}
		// Insere valores do formulário
		insForm.Exec(name, email)
		// exibe um log com os valores digitados np formulário
		log.Println("INSERT: Name: " + name + " | E-mail: " + email)
	}
	// encerra a conexão do dbConn()
	defer db.Close()

	// Retorna a Home
	http.Redirect(w, r, "/", 301)
}

// Função Update
func Update(w http.ResponseWriter, r *http.Request) {
	// Abre a conexão com o banco de dados usanndo a função dbConn()
	db := dbConn()
	// Verifica o METHOD do formulário passado
	if r.Method == "POST" {
		// Pega os campos do formulário
		name := r.FormValue("name")
		email := r.FormValue("email")
		id := r.FormValue("id")

		// Prepara o SQL e verifica erros
		insForm, err := db.Prepare("UPDATE names SET name=?, email=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		// Insere valores do formulário com a SQL tratada e verifica erros
		insForm.Exec(name, email, id)
		// Exibe um log com os valores digitados no formulario
		log.Println("UPDATE: Name: " + name + " |E-mail: " + email)
	}
	// Encerra a conexão
	defer db.Close()

	// Retorna a Home
	http.Redirect(w, r, "/", 301)
}

// Função Delete
func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	delForm, err := db.Prepare("DELETE FROM names WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	delForm.Exec(nId)
	log.Println("DELETE")
	defer db.Close()

	http.Redirect(w, r, "/", 301)
}
