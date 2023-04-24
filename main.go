package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"time"
)

type Player struct {
	Name  string
	Stats Stats
}

type Stats struct {
	HP  int
	MP  int
	STR int
	INT int
	LUC int
}

func main() {
	// 讀取兩個圖片檔案
	img1 := loadImage("image1.png")
	img2 := loadImage("image2.png")

	// 計算初五種數值
	stats1 := calculateStats(img1)
	stats2 := calculateStats(img2)

	// 將初五種數值作為兩個玩家的屬性
	player1 := Player{Name: "Player 1", Stats: stats1}
	player2 := Player{Name: "Player 2", Stats: stats2}

	// 開始對戰
	winner := fight(player1, player2)

	// 輸出勝利者
	fmt.Printf("%s wins!\n", winner.Name)
}

func loadImage(filename string) image.Image {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	return img
}

func calculateStats(img image.Image) Stats {
	stats := Stats{}

	// 計算圖片的寬度和高度
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 計算每個 pixel 的值並加總
	var r, g, b, a uint32
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			r, g, b, a = pixel.RGBA()

			stats.HP += int(r)
			stats.MP += int(g)
			stats.STR += int(b)
			stats.INT += int(a)
			stats.LUC += int(r + g + b + a)
		}
	}

	// 將總和除以 pixel 總數，得到平均值
	totalPixels := width * height
	stats.HP /= totalPixels
	stats.MP /= totalPixels
	stats.STR /= totalPixels
	stats.INT /= totalPixels
	stats.LUC /= totalPixels

	return stats
}

func fight(player1 Player, player2 Player) Player {
	fmt.Printf("%s vs. %s!\n", player1.Name, player2.Name)

	// 隨機決定哪個玩家先攻
	rand.Seed(time.Now().UnixNano())
	players := []Player{player1, player2}
	attackerIndex := rand.Intn(2)
	attacker := players[attackerIndex]
	defender := players[(attackerIndex+1)%2]
	fmt.Printf("%s attacks first!\n", attacker.Name)

	// 玩家攻
	for {
		// 計算攻擊傷害
		damage := attacker.Stats.STR + rand.Intn(attacker.Stats.INT)
		fmt.Printf("%s attacks for %d damage!\n", attacker.Name, damage)

		// 扣除血量
		defender.Stats.HP -= damage
		if defender.Stats.HP <= 0 {
			return attacker
		}

		// 切換攻擊者和防禦者
		attacker, defender = defender, attacker
	}
}
