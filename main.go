package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main(){
	handler:=http.NewServeMux()

	// C R U D
	//handler.HandleFunc("/hello/",Logger(helloHandler))

	handler.HandleFunc("/book/", Logger(bookHandler))

	handler.HandleFunc("/books/", Logger(booksHandler))

	s:= http.Server{
		Addr: ":8080",
		Handler: handler,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,        //1*2^20 - 128kBytes
	}
	log.Fatal(s.ListenAndServe())
}


func Logger(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		fmt.Printf("server [net/http] method [%s] connection from [%v]", r.Method, r.RemoteAddr)

		next.ServeHTTP(w,r)
	}
}


func bookHandler (w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet{
		fmt.Printf("Get")
		handleGetBook(w,r)
	}else if r.Method == http.MethodPost{
		fmt.Printf("POST")
		handleAddBook(w,r)
	}else if r.Method == http.MethodDelete{
		fmt.Printf("Delete")
		handleDeleteBook(w,r)
	}else if r.Method == http.MethodPut{
		fmt.Printf("Put")
		handleUpdateBook(w,r)
	}
}
func handleDeleteBook( w http.ResponseWriter, r *http.Request){

	id:= strings.Replace(r.URL.Path, "/book/", "",1)

	var resp Resp


	err := bookStore.DeleteBook(id)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)

		resp.Error = err.Error()

		respJson, _ := json.Marshal(resp)
		fmt.Printf("Error2")
		w.Write(respJson)
		return
	}
	booksHandler(w, r)
}

func handleUpdateBook( w http.ResponseWriter, r *http.Request){

	id:= strings.Replace(r.URL.Path, "/book/", "",1)
	decoder := json.NewDecoder(r.Body)

	var book Book
	var resp Resp

	err := decoder.Decode(&book)


	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()

		respJson, _ := json.Marshal(resp)
		fmt.Printf("Error1")
		w.Write(respJson)
		return
	}
	book.Id = id
	err = bookStore.UpdateBook(book)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)

		resp.Error = err.Error()

		respJson, _ := json.Marshal(resp)
		fmt.Printf("Error2")
		w.Write(respJson)
		return
	}
	resp.Message = book
	respJson, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)

}
func handleAddBook( w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	fmt.Printf("handleAddBook")
	var book Book
	var resp Resp //ответ

	err := decoder.Decode(&book)


	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()

		respJson, _ := json.Marshal(resp)
		fmt.Printf("Error1")
		w.Write(respJson)
		return
	}
	err = bookStore.AddBooks(book)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()

		respJson, _ := json.Marshal(resp)
		fmt.Printf("Error2")
		w.Write(respJson)
		return
	}

	booksHandler(w, r)

}

func booksHandler (w http.ResponseWriter, r *http.Request){ //вывод на экран
	//if r.Method == http.MethodGet{
	//	handleGetBook(w,r)
//	}
    fmt.Printf("booksHandler - vyvod")
	w.WriteHeader(http.StatusOK)

	resp := Resp{
		 Message: bookStore.GetBooks(),
	}

	bookJson, _ := json.Marshal(resp)

	w.Write(bookJson)
}


func handleGetBook (w http.ResponseWriter, r *http.Request){
	fmt.Printf("handleGetBook")

	id:= strings.Replace(r.URL.Path, "/book/", "",1)
	book := bookStore.FindBookByID(id)

	var resp Resp

	if book == nil{
		w.WriteHeader(http.StatusNotFound)
		resp.Error = fmt.Sprintf("")

		respJson, _ := json.Marshal(resp)
		fmt.Printf("Error3")
		w.Write(respJson)
		return
	}
	resp.Message = book
	respJson, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)

	w.Write(respJson)



}

type Resp struct{
	Message interface{}
	Error string
}

type Book struct{
	Id string `json:id`
	Author string `json:"author"`
	Name string `json:"name"`

}

type BookStore struct{
	books []Book
}

var bookStore = BookStore{
	books: make([]Book,0),
}

func (s BookStore) FindBookByID(id string) *Book {
	for _, book := range s.books{
		if book.Id == id{
			return &book
		}
	}
	return  nil
}

func (s BookStore) GetBooks() []Book {
	return s.books
}


func (s *BookStore) AddBooks(book Book) error {
	for _, bk := range s.books{
		if bk.Id == book.Id{
			return errors.New(fmt.Sprintf("Book with id %s not found", book.Id))
		}
	}
	s.books = append(s.books,book)
	return nil
}

func (s *BookStore) UpdateBook(book Book) error {
	for i, bk := range s.books{
		if bk.Id == book.Id{
			s.books[i] = book
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Book with id %s not found", book.Id))
}

func (s *BookStore) DeleteBook(id string) error{
	for i,bk := range s.books{
		if bk.Id == id{
			s.books = append(s.books[:i], s.books[i+1:]...)

			return nil
		}
	}
	return errors.New(fmt.Sprintf("Book with id %s not found", id))
}
/*
func helloHandler(w http.ResponseWriter, r *http.Request) {


	name:= strings.Replace(r.URL.Path, "/hello/", "",1)

	resp := Resp {
		Message: fmt.Sprintf("hello %s. Glad to see you again", name),
	}

	respJson, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)

	w.Write(respJson)

}*/
