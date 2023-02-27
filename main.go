package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	quiz := readQuiz("sources/quizzes.txt")
	hint := readQuiz("sources/hint.txt")
	genIndex(quiz, hint)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "output.html", gin.H{
			"": "",
		})
	})

	r.POST("/", func(c *gin.Context) {
		counter := 0

		form := []string{}
		quizzes := readQuiz("sources/quizzes.txt")
		answer := readQuiz("sources/answer.txt")

		//根据问题生成提交表单数量
		for key := range quizzes {
			form = append(form, c.PostForm("q"+strconv.Itoa(key+1)))
		}

		//判断是否正确
		for key, value := range form {
			if value == answer[key] {
				counter++
			}
		}

		//可调整做对的门限
		if counter <= 10 {
			c.JSON(200, gin.H{"error": "做对了" + strconv.Itoa(counter) + "道题目"})
		} else {
			c.JSON(200, gin.H{"congraduation!": "flag{test_flag}"})
		}
	})

	// 启动Web服务器
	r.Run(":8888")
}

// 打开问题文件，逐行读取
func readQuiz(fileName string) (quiz []string) {
	file, err := os.Open(fileName) // 打开文件
	if err != nil {
		panic(err)
	}
	defer file.Close() // 关闭文件

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		quiz = append(quiz, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return quiz
}

func genIndex(quizzes []string, hint []string) {
	quizzHtml := ""
	part1 := `<div class="card"><div class="card-body question-main"><div class="question-q">`
	part2 := `<br><small class="text-muted" font-size="5">`
	part3 := `</small></div><div class="input-group input-group-sm mb-3"><div class="input-group-prepend"><span class="input-group-text" id="q%s">clue</span></div><input type="text" name="q%s" class="form-control" aria-describedby="q%s" value=""></div></div></div><br>`

	for key, value := range quizzes {
		quizzHtml += part1 + value + part2 + hint[key] + fmt.Sprintf(part3, strconv.Itoa(key+1), strconv.Itoa(key+1), strconv.Itoa(key+1))
	}

	err := ioutil.WriteFile("temp/index.p2", []byte(quizzHtml), 0644)
	if err != nil {
		panic(err)
	}

	files := []string{"temp/index.p1", "temp/index.p2", "temp/index.p3"}
	output := "templates/output.html"

	out, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	for _, file := range files {
		in, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			panic(err)
		}
	}
}
