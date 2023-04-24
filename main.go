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

func NewPlayer(name string, stats Stats) *Player {
	return &Player{
		Name:  name,
		Stats: stats,
	}
}

func main() {
	// 讀取兩個圖片檔案，解析成人物資料
	player1Data, err := readImage("player1.png")
	if err != nil {
		fmt.Println("Error reading player1.png:", err)
		return
	}
	player2Data, err := readImage("player2.png")
	if err != nil {
		fmt.Println("Error reading player2.png:", err)
		return
	}

	// 根據解析出來的資料，創建兩個人物
	player1 := NewPlayer(player1Data.Name, player1Data.Stats)
	player2 := NewPlayer(player2Data.Name, player2Data.Stats)

	// 讓這兩個人物進行對戰
	winner := fight(player1, player2)

	// 印出勝利者的名字
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

func readImage(filePath string) (*Player, error) {
	// 開啟圖片檔案
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 讀取圖片檔案
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// 創建一個新的 CharacterData 物件
	data := &Player{}

	// 解析像素資訊，設定人物名字和屬性值
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			switch {
			case x < bounds.Max.X/2 && y < bounds.Max.Y/2:
				data.Name += string(r>>8) + string(g>>8) + string(b>>8)
			case x >= bounds.Max.X/2 && y < bounds.Max.Y/2:
				data.Stats.HP += int(r>>8) + int(g>>8) + int(b>>8)
			case x < bounds.Max.X/2 && y >= bounds.Max.Y/2:
				data.Stats.MP += int(r>>8) + int(g>>8) + int(b>>8)
			case x >= bounds.Max.X/2 && y >= bounds.Max.Y/2:
				data.Stats.STR += int(r >> 8)
				data.Stats.INT += int(g >> 8)
				data.Stats.LUC += int(b >> 8)
			}
		}
	}

	return data, nil
}

func fight(player1 *Player, player2 *Player) Player {
	fmt.Printf("%s vs. %s!\n", player1.Name, player2.Name)

	// 隨機決定哪個玩家先攻
	rand.Seed(time.Now().UnixNano())
	players := []Player{*player1, *player2}
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
