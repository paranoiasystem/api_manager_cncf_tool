package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		//db, err := database.OpenConnection()
		//if err != nil {
		//	log.Fatal("failed to connect database: %v", err)
		//}
		//defer db.Close()
		//
		//rows, err := db.Query("SELECT version()")
		//if err != nil {
		//	log.Fatalf("query failed: %v", err)
		//}
		//defer rows.Close()
		//
		//// Scansiona il risultato della query
		//var version string
		//for rows.Next() {
		//	if err := rows.Scan(&version); err != nil {
		//		log.Fatalf("failed to scan result: %v", err)
		//	}
		//}
		//
		//c.JSON(200, gin.H{
		//	"db_version": version,
		//})
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
}
