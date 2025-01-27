package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

// Connect to PostgreSQL
func init() {
	var err error
	connStr := "user=user password=password dbname=ahnenbaum_db sslmode=disable" // Replace with your DB credentials
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
}

func main() {
	router := gin.Default()

	// Family Descendant Routes
	router.GET("/api/descendants", getDescendants)
	router.GET("/api/descendants/:id", getDescendantByID)
	router.POST("/api/descendants", createDescendant)
	router.PUT("/api/descendants/:id", updateDescendant)
	router.DELETE("/api/descendants/:id", deleteDescendant)

	// Attachment Figures Routes
	router.GET("/api/attachment_figures", getAttachmentFigures)
	router.POST("/api/attachment_figures", createAttachmentFigures)

	// attachment joins
	router.POST("/api/attachmentjoin", createAttachmentjoin)
	router.PUT("/api/attachmentjoins/:id", updateAttachmentJoin)
	router.DELETE("/api/attachmentjoins/:id", deleteAttachmentJoin)
	// Start the server
	router.Run(":8080")
}

// --- CRUD Functions for family_descendant ---

// Get all descendants
func getDescendants(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, description, image, generation, gender, birth_date, death_date FROM family_descendant")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var descendants []map[string]interface{}
	for rows.Next() {
		var id, generation int
		var name, description, image, gender string
		var birthDate, deathDate sql.NullString

		if err := rows.Scan(&id, &name, &description, &image, &generation, &gender, &birthDate, &deathDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		descendants = append(descendants, gin.H{
			"id":          id,
			"name":        name,
			"description": description,
			"image":       image,
			"generation":  generation,
			"gender":      gender,
			"birth_date":  birthDate.String,
			"death_date":  deathDate.String,
		})
	}

	c.JSON(http.StatusOK, descendants)
}

// Get a single descendant by ID
func getDescendantByID(c *gin.Context) {
	id := c.Param("id")
	row := db.QueryRow("SELECT id, name, description, image, generation, gender, birth_date, death_date FROM family_descendant WHERE id = $1", id)

	var descendant map[string]interface{}
	var idInt, generation int
	var name, description, image, gender string
	var birthDate, deathDate sql.NullString

	if err := row.Scan(&idInt, &name, &description, &image, &generation, &gender, &birthDate, &deathDate); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Descendant not found"})
		return
	}

	descendant = gin.H{
		"id":          idInt,
		"name":        name,
		"description": description,
		"image":       image,
		"generation":  generation,
		"gender":      gender,
		"birth_date":  birthDate.String,
		"death_date":  deathDate.String,
	}

	c.JSON(http.StatusOK, descendant)
}

// Create a new descendant
func createDescendant(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Image       string `json:"image"`
		Generation  int    `json:"generation"`
		Gender      string `json:"gender"`
		BirthDate   string `json:"birth_date"`
		DeathDate   string `json:"death_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO family_descendant (name, description, image, generation, gender, birth_date, death_date) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var id int
	err := db.QueryRow(query, input.Name, input.Description, input.Image, input.Generation, input.Gender, input.BirthDate, input.DeathDate).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// Update an existing descendant
func updateDescendant(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
		Generation  int    `json:"generation"`
		Gender      string `json:"gender"`
		BirthDate   string `json:"birth_date"`
		DeathDate   string `json:"death_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE family_descendant 
              SET name = $1, description = $2, image = $3, generation = $4, gender = $5, birth_date = $6, death_date = $7 
              WHERE id = $8`
	_, err := db.Exec(query, input.Name, input.Description, input.Image, input.Generation, input.Gender, input.BirthDate, input.DeathDate, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Descendant updated successfully"})
}

// Delete a descendant
func deleteDescendant(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM family_descendant WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Descendant deleted successfully"})
}

// --- CRUD Functions for attachment_figures ---

// Get all attachment figures
func getAttachmentFigures(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, description, image, birth_date, death_date, gender FROM attachment_figures")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var attachmentfigures []map[string]interface{}
	for rows.Next() {
		var id int
		var name, description, image, gender string
		var birthDate, deathDate sql.NullString

		if err := rows.Scan(&id, &name, &description, &image, &birthDate, &deathDate, &gender); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		attachmentfigures = append(attachmentfigures, gin.H{
			"id":          id,
			"name":        name,
			"description": description,
			"image":       image,
			"birth_date":  birthDate.String,
			"death_date":  deathDate.String,
			"gender":      gender,
		})
	}

	c.JSON(http.StatusOK, attachmentfigures)
}

// Create a new attachment
func createAttachmentFigures(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Image       string `json:"image"`
		Gender      string `json:"gender"`
		BirthDate   string `json:"birth_date"`
		DeathDate   string `json:"death_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO attachment_figures (name, description, image, gender, birth_date, death_date) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int
	err := db.QueryRow(query, input.Name, input.Description, input.Image, input.Gender, input.BirthDate, input.DeathDate).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// CRUD for attachment figures
func createAttachmentjoin(c *gin.Context) {
	var input struct {
		attachment_figure_id int    `json:"attachment_figure_id" binding:"required"`
		descendant_id        string `json:"descendant_id" binding:"required"`
	}
	query := `INSERT INTO descendant_attachments (attachment_figure_id, descendant_id) 
              VALUES ($1, $2) RETURNING id`
	var id int
	err := db.QueryRow(query, input.attachment_figure_id, input.descendant_id).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func updateAttachmentJoin(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		AttachmentFigureID int `json:"attachment_figure_id" binding:"required"`
		DescendantID       int `json:"descendant_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE descendant_attachments 
              SET attachment_figure_id = $1, descendant_id = $2 
              WHERE id = $3`
	_, err := db.Exec(query, input.AttachmentFigureID, input.DescendantID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attachment join updated successfully"})
}

func deleteAttachmentJoin(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM descendant_attachments WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attachment join deleted successfully"})
}

func getAllDescendants(c *gin.Context) {
	id := c.Param("id") // ID of the individual whose descendants to find

	query := `
		WITH RECURSIVE descendants AS (
			SELECT id, name, family_parent, related_by_attachment, generation
			FROM family_descendant
			WHERE id = $1
			UNION ALL
			SELECT fd.id, fd.name, fd.family_parent, fd.related_by_attachment, fd.generation
			FROM family_descendant fd
			INNER JOIN descendants d ON fd.family_parent = d.id
		)
		SELECT id, name, family_parent, related_by_attachment, generation FROM descendants;
	`

	rows, err := db.Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, familyParent, relatedByAttachment, generation int
		var name string

		if err := rows.Scan(&id, &name, &familyParent, &relatedByAttachment, &generation); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result = append(result, gin.H{
			"id":                    id,
			"name":                  name,
			"family_parent":         familyParent,
			"related_by_attachment": relatedByAttachment,
			"generation":            generation,
		})
	}

	c.JSON(http.StatusOK, result)
}

func getAllParentsAndAttachedFigures(c *gin.Context) {
	id := c.Param("id") // ID of the individual whose parents and attached figures to find

	query := `
		SELECT 
			fd.id AS descendant_id, 
			fp.id AS family_parent_id, 
			fp.name AS family_parent_name, 
			af.id AS attachment_figure_id, 
			af.name AS attachment_figure_name
		FROM family_descendant fd
		LEFT JOIN family_descendant fp ON fd.family_parent = fp.id
		LEFT JOIN attachment_figures af ON fd.related_by_attachment = af.id
		WHERE fd.id = $1;
	`

	row := db.QueryRow(query, id)

	var descendantID, familyParentID, attachmentFigureID sql.NullInt32
	var familyParentName, attachmentFigureName sql.NullString

	if err := row.Scan(&descendantID, &familyParentID, &familyParentName, &attachmentFigureID, &attachmentFigureName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found for the given ID"})
		return
	}

	result := gin.H{
		"descendant_id":     descendantID.Int32,
		"family_parent":     gin.H{"id": familyParentID.Int32, "name": familyParentName.String},
		"attachment_figure": gin.H{"id": attachmentFigureID.Int32, "name": attachmentFigureName.String},
	}

	c.JSON(http.StatusOK, result)
}
