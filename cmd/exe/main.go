package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(http.StatusOK, struct{}{})
}

func ping(c *gin.Context) {
	cmd := exec.Command("mplayer", "-fs", "-vo", "xv", "-ao", "alsa:device=hdmi", "/home/mario/Videos/Nuovo Cinema Paradiso/Nuovo.cinema.Paradiso.(1988).BDRip.720p.AC3.X264-CHD-Italian.mkv", "&")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

	c.JSON(200, gin.H{
		"message": "pong",
	})

	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s\n", out)

	//mplayer -fs -vo xv -ao alsa:device=hdmi /home/mario/Videos/Nuovo\ Cinema\ Paradiso/Nuovo.cinema.Paradiso.\(1988\).BDRip.720p.AC3.X264-CHD-Italian.mkv
}

func main() {
	g := gin.Default()

	g.GET("/ping", ping)
	g.OPTIONS("/ping", preflight)

	g.Run(":5000")
}
