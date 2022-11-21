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

func main() {
	outFile()
	outCsv()
	outStd()
	http.HandleFunc("/", outJSON)
	http.ListenAndServe(":8080", nil)
}