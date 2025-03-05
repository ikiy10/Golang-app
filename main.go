package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// Struktur data untuk tugas (To-Do)
type Todo struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status string `json:"status"`
}

// Data sementara (database tiruan)
var todos = []Todo{
    {ID: 1, Title: "Belajar Golang", Status: "pending"},
    {ID: 2, Title: "Bikin API dengan Gin", Status: "done"},
}

// GET: Ambil semua tugas
func getTodos(c *gin.Context) {
    c.JSON(http.StatusOK, todos)
}

// POST: Tambah tugas baru (single atau multiple)
func addTodos(c *gin.Context) {
    var newTodos []Todo
    var singleTodo Todo

    // Coba bind JSON ke dalam slice (array)
    if err := c.ShouldBindJSON(&newTodos); err != nil {
        // Jika gagal, coba bind JSON ke dalam satu objek
        if err := c.ShouldBindJSON(&singleTodo); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
            return
        }
        // Jika berhasil, ubah menjadi array
        newTodos = []Todo{singleTodo}
    }

    // Cari ID terakhir agar tidak ada ID yang duplikat
    lastID := 0
    if len(todos) > 0 {
        lastID = todos[len(todos)-1].ID
    }

    // Tambahkan data baru ke daftar todos
    for i := range newTodos {
        lastID++
        newTodos[i].ID = lastID
        todos = append(todos, newTodos[i])
    }

    c.JSON(http.StatusCreated, newTodos)
}


// PUT: Perbarui status tugas
func updateTodo(c *gin.Context) {
    var updatedTodo Todo
    if err := c.ShouldBindJSON(&updatedTodo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for i, t := range todos {
        if t.ID == updatedTodo.ID {
            todos[i].Title = updatedTodo.Title
            todos[i].Status = updatedTodo.Status
            c.JSON(http.StatusOK, todos[i])
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"message": "Tugas tidak ditemukan"})
}

// DELETE: Hapus tugas berdasarkan ID
func deleteTodo(c *gin.Context) {
    var todoID struct {
        ID int `json:"id"`
    }
    if err := c.ShouldBindJSON(&todoID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for i, t := range todos {
        if t.ID == todoID.ID {
            todos = append(todos[:i], todos[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "Tugas berhasil dihapus"})
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"message": "Tugas tidak ditemukan"})
}

// Main function untuk menjalankan server
func main() {
    r := gin.Default()

    // Routing
    r.GET("/todos", getTodos)       // Lihat semua tugas
    r.POST("/todos", addTodos)      // Tambah satu atau banyak tugas
    r.PUT("/todos", updateTodo)     // Perbarui tugas
    r.DELETE("/todos", deleteTodo)  // Hapus tugas

    r.Run(":8080") // Jalankan server di port 8080
}
