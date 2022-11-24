package main

import (
	"fmt"
	"os"
	"log"
	"encoding/csv"
	"encoding/json"
	"compress/gzip"
	"io"
	"net/http"

	"bufio"
	"strings"
	"bytes"

	"crypto/rand"

)
/* 問題演習
①ファイルに対するフォーマット出力
②CSV出力
③gzipされたJSON出力しながら標準出力にログ出力
*/

//①
func outFile() {
	file, err := os.Create("test.txt") //ファイル新規作成
	if(err != nil) {
		panic(err)
	}
	//Fprintfでio.Writerに数値や文字列を出力できる
	fmt.Fprintf(file, "%d\n", 5)
	fmt.Fprintf(file, "%s\n", 5) //%!s(int=5)
	fmt.Fprintf(file, "%d\n", 5.3) //%!d(float64=5.3)

	fmt.Fprintf(file, "%d\n", 5)
	fmt.Fprintf(file, "%f\n", 5.3) //%!s(int=5)
	fmt.Fprintf(file, "%v\n", 5)

	fmt.Fprintf(file, "%d\n", "string") //%!d(string=string)
	fmt.Fprintf(file, "%s\n", "string")
	fmt.Fprintf(file, "%v\n", "string")
}

//②
func outCsv() { //csvファイルに出力
	file, err := os.Create("practice.csv")
	if (err != nil) {
		panic(err)
	}

	records := [][]string{
		[]string{"名前", "年齢", "出身地", "性別"},
		[]string{"山本", "24", "兵庫", "男性"},
		[]string{"鈴木", "25", "大阪", "女性"},
		[]string{"鎌田", "27", "東京", "男性"},
	}

	defer file.Close() //closeしないと書き込みが終われない

	w := csv.NewWriter(file)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	defer w.Flush() //書かないと出力されない　バッファに 残っているデータをすべて書き込む
}

func outStd() { //標準出力ver
	records := [][]string{
		[]string{"名前", "年齢", "出身地", "性別"},
		[]string{"山本", "24", "兵庫", "男性"},
		[]string{"鈴木", "25", "大阪", "女性"},
		[]string{"鎌田", "27", "東京", "男性"},
	}

	w := csv.NewWriter(os.Stdout)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()
}

//③
/*
1. JSON作成
2. os.StdOutに出力するとログが出るようにJSONを文字列に変換
3. gzip圧縮を行いながら圧縮前の出力を標準出力にも出すように io.MultiWriterを使う
4. gzip出力の最後にはFlush()必要です。
*/
func outJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")

	sourse := map[string]string{ //	//json化する元データ
		"Hello": "World",
	}
	//ここにコードを書く
	jSourse, _ := json.Marshal(sourse) //json化 　_ なし：assignment mismatch: 1 variable but json.Marshal returns 2 values
	newJsourse := string(jSourse) //jsonの文字列化
	// fmt.Printf("%T\n", newJsourse) //string

	//jsonをzipに圧縮
	file, err := os.Create("goPractice.txt.gz")
	if err != nil {
		log.Fatal(err)
	}
	writer := gzip.NewWriter(file)
	writer.Header.Name = "goPractice.txt"

	logW := io.MultiWriter(os.Stdout, writer) //標準出力とzipにjson出力
	io.WriteString(logW, /*jSourse*/newJsourse)
	defer writer.Close()
	defer writer.Flush()

}


//io.Reader練習
//例題
var sourse = `1行目
2行目
3行目`

//改行で区切る
func nReader() {
	reader := bufio.NewReader(strings.NewReader(sourse))
	for {
		line, err := reader.ReadString('\n')
		fmt.Printf("%#v\n", line)
		if err == io.EOF {
			break
		}
	}
}

//データ型を指定して解析
func dataReader() {
	var sourse = "123 1.234 1.0e4 test"
	reader := strings.NewReader(sourse)
	var i int
	var f, g float64
	var s string
	fmt.Fscan(reader, &i, &f, &g, &s) //fmt.Fscan(reader)がデータがスペース区切りである全体
	fmt.Printf("i=%#v\n f=%#v\n g=%#v\n s=%#v\n", i, f, g, s)
}

//csv形式を解析
func csvReader() {
	var csvSourse = 
	`13101,"100  ","1000003","ﾄｳｷｮｳﾄ","ﾁﾖﾀﾞｸ","ﾋﾄﾂﾊﾞｼ"(1ﾁｮｳﾒ)","東京都","千代田区","一ツ橋（一丁目）",1,0,1,0,0,0
	13101,"101  ","1010003","ﾄｳｷｮｳﾄ","ﾁﾖﾀﾞｸ","ﾋﾄﾂﾊﾞｼ"(2ﾁｮｳﾒ)","東京都","千代田区","一ツ橋（二丁目）",1,0,1,0,0,0
	13101,"100  ","1000012","ﾄｳｷｮｳﾄ","ﾁﾖﾀﾞｸ","ﾋﾋﾞﾔｺｳｴﾝ","東京都","千代田区","日比谷公園",0,0,0,0,0,0
	13101,"102  ","1000093","ﾄｳｷｮｳﾄ","ﾁﾖﾀﾞｸ","ﾋﾗｶﾜﾁｮｳ","東京都","千代田区","平河町",0,0,1,0,0,0
	13101,"102  ","1000071","ﾄｳｷｮｳﾄ","ﾁﾖﾀﾞｸ","ﾌｼﾞﾐ","東京都","千代田区","富士見",0,0,1,0,0,0`

	reader := strings.NewReader(csvSourse)
	csvReader := csv.NewReader(reader)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(line[2], line[3:5])
	}
}

//io.Readerの全ての入力がつながっているように動作
func multiReader() {
	header := bytes.NewBufferString("----HEADER----\n")
	content := bytes.NewBufferString("Example of io.MultiReader\n")
	footer := bytes.NewBufferString("----FOOTER----\n")

	reader := io.MultiReader(header, content, footer)
	io.Copy(os.Stdout, reader)
}

//問題

/*
1.ファイルのコピー
①ファイル準備
②読み込み
③コピー
*/
func oldNew() {
	file, err := os.Create("old.txt")
	if err != nil {
		log.Fatalln("create old file error:", err)
	}

	// if err := file.Write("context in old file"); err != nil { //byte文字列のみ
	// 	log.Fatalln("error writing to old:", err)
	// }
	io.WriteString(file, "context in old file\n")

	// io.Copy(os.Stdout, file) //出力されない
	//参考	logW := io.MultiWriter(os.Stdout, writer) //標準出力とzipにjson出力

	newFile, err := os.Create("new.txt")
	if err != nil {
		log.Fatalln("create new file error:", err)
	}
	
	openFile, err := os.Open("old.txt")
	if err != nil {
		log.Fatalln("open file error:", err)
	}

	defer openFile.Close()
	io.Copy(newFile, openFile)
	
}

//2.テスト用の適当なサイズのファイルを作成。ファイルを作成してランダムな内容で埋める
/*
1 1024バイトの長さのバイナリファイル作成
2 用意したファイルに書き込み
👇
1024バイトのランダムな文字列作成⇒書き込む
*/
func randFile() {
	buffer := make([]byte, 1024) //バッファを準備してる
	rand.Read(buffer)//指定した長さのバイト配列を生成

	testFile, err := os.Create("test.txt")
	if err != nil {
		log.Fatalln("create file error!", err)
	}

	//openFile := os.openFile("test.txt")//os.openFile(testFile)　※()内はファイル名
	defer testFile/*openFile*/.Close()

	io.WriteString(testFile/*openFile*/, string(buffer))//io.WriteString(openFile, buffer)
}


func main() {
	outFile()
	outCsv()
	outStd()
	nReader()
	dataReader()
	csvReader()
	multiReader()
	oldNew()
	randFile()
	http.HandleFunc("/", outJSON)
	http.ListenAndServe(":8080", nil)
}