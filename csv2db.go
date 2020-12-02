//Đọc file csv và tạo ra bảng coutry
package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

//Hàm kiểm lỗi
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

//Kết nối với database
func connect(host, user, password, databaseName, port string) *sql.DB {
	msql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, databaseName)
	db, err := sql.Open(databaseName, msql)
	checkError(err)
	return db
}

//Tạo ra bảng Country trên database
func createCountryTable(db *sql.DB) {
	_, e := db.Exec("CREATE TABLE Country (Country_code	varchar(10) PRIMARY KEY, Country_Name varchar(50))")
	checkError(e)
}

//Thêm 1 phần tử vào bảng
func insertTable(db *sql.DB, countryCode string, countryName string) {
	_, e := db.Exec("INSERT INTO country(Country_code, Country_Name) values($1,$2)", countryCode, countryName)
	checkError(e)
}

//Đọc file csv
func readCsv(filename string) [][]string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	record, err := csv.NewReader(f).ReadAll()
	checkError(err)
	return record
}

//Thêm các phần tử vào table
func insertValue(db *sql.DB, record [][]string) {
	flag := false
	code := -1
	name := -1
	for _, value := range record {
		if !flag {
			if value[0] == "Name" {
				name = 0
				code = 1
			} else {
				name = 1
				code = 0
			}
			flag = true
		} else {
			insertTable(db, value[code], value[name])
		}
	}
}

//Trợ giúp
func help() {
	fmt.Println("To run this file, go to the location of cvs2bs.exe, on cmd type csv2db.exe args\n")
	fmt.Println("For example\n")
	fmt.Println("D:/cvs2bd.exe E:/country.csv localhost 1234 nameDatabase 5678 admin adminPassword\n")
	fmt.Println("E:/country_csv is the place you put the csv file")
	fmt.Println("localhost is your hostName")
	fmt.Println("1234 is port of host")
	fmt.Println("nameDatabase is your databaseName")
	fmt.Println("admin is the user of database")
	fmt.Println("adminPassword is the password of admin\n")
	os.Exit(0)
}

//Kiểm tra các tham số nhập vào
func checkArgs() {
	if len(os.Args) <= 1 {
		fmt.Printf("/? for help")
		os.Exit(0)
	}
	if os.Args[1] == "/?" {
		help()
	}
	if len(os.Args) < 7 {
		fmt.Println("Missing some args. /? for more infomation")
		os.Exit(0)
	}
}

func main() {
	checkArgs()
	csvpath := os.Args[1]
	host := os.Args[2]
	port := os.Args[3]
	databaseName := os.Args[4]
	user := os.Args[5]
	password := os.Args[6]

	record := readCsv(csvpath)
	db := connect(host, user, password, databaseName, port)
	createCountryTable(db)
	insertValue(db, record)
	fmt.Println("Create Table Success")
	defer db.Close()
	fmt.Scanln()
}
