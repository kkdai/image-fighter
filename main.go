package main

import (
	"fmt"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
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
	player1Data, err := readImage("./player1.jpg")
	if err != nil {
		fmt.Println("Error reading player1.jpg:", err)
		return
	}
	player2Data, err := readImage("./player2.jpg")
	if err != nil {
		fmt.Println("Error reading player2.png:", err)
		return
	}

	// 根據解析出來的資料，創建兩個人物
	player1 := NewPlayer(player1Data.Name, player1Data.Stats)
	fmt.Println("Player1 status:", player1)
	player2 := NewPlayer(player2Data.Name, player2Data.Stats)
	fmt.Println("Player2 status:", player2)

	// 讓這兩個人物進行對戰
	winner := fight(player1, player2)

	// 印出勝利者的名字
	fmt.Printf("%s wins!\n", winner.Name)
}

func readImage(filePath string) (*Player, error) {
	// 開啟圖片檔案
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 讀取圖片檔案
	img, err := jpeg.Decode(file)
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
				// 隨機生成姓氏和名字
				data.Name = filePath
			case x >= bounds.Max.X/2 && y < bounds.Max.Y/2:
				data.Stats.HP += int(r>>8) + int(g>>8) + int(b>>8)
				if data.Stats.HP > 65535 {
					data.Stats.HP = 65535
				}
			case x < bounds.Max.X/2 && y >= bounds.Max.Y/2:
				data.Stats.MP += int(r>>8) + int(g>>8) + int(b>>8)
				if data.Stats.MP > 65535 {
					data.Stats.MP = 65535
				}
			case x >= bounds.Max.X/2 && y >= bounds.Max.Y/2:
				data.Stats.STR += int(r >> 8)
				if data.Stats.STR > 255 {
					data.Stats.STR = 255
				}
				data.Stats.INT += int(g >> 8)
				if data.Stats.INT > 255 {
					data.Stats.INT = 255
				}
				data.Stats.LUC += int(b >> 8)
				if data.Stats.LUC > 255 {
					data.Stats.LUC = 255
				}
			}
		}
	}

	return data, nil
}

// fight: 在這裡我們先比較兩個玩家的 LUC 屬性值，較高的玩家被選為先攻。
// 如果兩個玩家的 LUC 屬性值相同，則隨機選擇一個玩家作為先攻。這樣就可以確保 LUC 更高的玩家先攻了。
func fight(player1 *Player, player2 *Player) Player {
	fmt.Printf("%s vs. %s!\n", player1.Name, player2.Name)

	// 選出 LUC 更高的玩家先攻
	var attacker, defender Player
	if player1.Stats.LUC > player2.Stats.LUC {
		attacker = *player1
		defender = *player2
	} else {
		attacker = *player2
		defender = *player1
	}
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
