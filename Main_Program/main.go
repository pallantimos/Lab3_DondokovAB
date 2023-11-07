package Main_Program

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	_ "github.com/lib/pq"
)

type Registration struct {
}

func (r Registration) checkRegistrate(login string, pass string, pass2 string) (string, bool) {
	err := ""
	isCorrect := true

	err, isCorrect = r.checkPass(pass, pass2)

	if isCorrect {
		err, isCorrect = r.checkLogin(login)
	}

	return err, isCorrect
}

func (r Registration) checkLogin(login string) (string, bool) {
	listLogin := [5]string{"Aldar", "Aleksey", "Ivan", "Mikhail", "Krug"}

	for i := 0; i < len(listLogin); i++ {
		if login == listLogin[i] {
			return "Логин уже существует", false
		}

	}

	if login == "" {
		return "Пустая строка в качества логина", false
	}

	if utf8.RuneCountInString(login) < 5 {
		return "Логин меньше 5 символов", false
	}

	reLogin := regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(login)
	reMail := regexp.MustCompile(`^\w+@\w+\.\w+$`).MatchString(login)
	rePhone := regexp.MustCompile(`^\+\d{1,3}-\d{3}-\d{3}-\d{4}$`).MatchString(login)

	if containsPlus := strings.Contains(login, "+"); containsPlus && !rePhone {
		return "Номер телефона не удовлетворяет заданному формату +x-xxx-xxx-xxxx", false
	}

	if containsAt := strings.Contains(login, "@"); containsAt && !reMail {
		return "Email не удовлетворяет общему формату xxx@xxx.xxx", false
	}

	if !reLogin && !reMail && !rePhone {
		return "Логин содержит символы, отличные от латиницы, цифр и знака подчеркивания", false
	}

	return "", true
}

func (r Registration) checkPass(pass string, pass2 string) (string, bool) {
	isUpperLetter := false
	isDownLetter := false
	isDigit := false
	isSymbol := false

	if pass == "" || pass2 == "" {
		return "Пустая строка в качестве пароля", false
	}

	for _, r := range pass {
		if unicode.Is(unicode.Latin, r) {
			return "Пароль содержит латиницу", false
		}
		if unicode.IsUpper(r) {
			isUpperLetter = true
		} else if unicode.IsLetter(r) {
			isDownLetter = true
		} else if unicode.IsDigit(r) {
			isDigit = true
		} else {
			isSymbol = true
		}
	}

	if !isDownLetter {
		return "Пароль не содержит строчную букву", false
	}

	if !isUpperLetter {
		return "Пароль не содержит заглавную букву", false
	}

	if !isSymbol {
		return "Пароль не содержит спецсимвола", false
	}

	if !isDigit {
		return "Пароль не содержит цифру", false
	}

	if utf8.RuneCountInString(pass) < 7 {
		return "Пароль меньше 7 символов", false
	}

	if pass != pass2 {
		return "Пароли не совпадают", false
	}

	return "", true
}

type DbWorker struct {
}

func (db DbWorker) Add(login string, pass string, pass2 string) string {
	dbpg, err := sql.Open("postgres", "user=Aldar password=123 dbname=Lab3 sslmode=disable")
	if err != nil {
		panic(err)
	}

	r := new(Registration)

	var err_auth string
	var result_auth bool
	err_auth, result_auth = r.checkRegistrate(login, pass, pass2)

	_, err = dbpg.Exec(fmt.Sprintf("INSERT INTO public.users(login, pass, pass2, result_auth, err_msg)VALUES ('%s', '%s', '%s', '%t', '%s');", login, pass, pass2, result_auth, err_auth))
	if err != nil {
		panic(err)
	}

	rows, err := dbpg.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var login string
		var pass string
		var pass2 string
		var errorString string

		rows.Scan()

		err = rows.Scan(&login, &pass, &pass2, &result_auth, &errorString)
		if err != nil {
			panic(err)
		}
		print(login)
	}

	if err_auth != "" {
		return err_auth
	}

	defer dbpg.Close()
	// Ваш код для выполнения операции добавления записи в базу данных
	// Например, вы можете использовать dbConn.Exec() для выполнения SQL-запроса
	return ""
}

func (db DbWorker) SelectUsers() [][]string {
	dbpg, err := sql.Open("postgres", "user=Aldar password=123 dbname=Lab3 sslmode=disable")
	if err != nil {
		panic(err)
	}

	rows, err := dbpg.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	var user_arr [][]string

	for rows.Next() {
		user := make([]string, 5)
		err = rows.Scan(&user[0], &user[1], &user[2], &user[3], &user[4])
		if err != nil {
			panic(err)
		}
		user_arr = append(user_arr, user)
	}

	defer dbpg.Close()
	return user_arr
}

func (db DbWorker) SelectUser(login string) [5]string {
	dbpg, err := sql.Open("postgres", "user=Aldar password=123 dbname=Lab3 sslmode=disable")
	if err != nil {
		panic(err)
	}

	rows, err := dbpg.Query(fmt.Sprintf("SELECT * FROM users WHERE login = '%s'", login))
	if err != nil {
		panic(err)
	}

	var user_arr [5]string

	rows.Next()

	err = rows.Scan(&user_arr[0], &user_arr[1], &user_arr[2], &user_arr[3], &user_arr[4])
	if err != nil {
		panic(err)
	}

	defer dbpg.Close()

	return user_arr
}

func (db DbWorker) DeleteUser(login string) {
	dbpg, err := sql.Open("postgres", "user=Aldar password=123 dbname=Lab3 sslmode=disable")
	if err != nil {
		panic(err)
	}

	_, err = dbpg.Exec(fmt.Sprintf("DELETE FROM public.users WHERE login = '%s';", login))
	if err != nil {
		panic(err)
	}
	defer dbpg.Close()
}

type Conrtoller struct {
}

func (ct Conrtoller) GetInfo() {
	fmt.Println("Введите что хотите сделать 1 - запись, 2 - удаление, 3 - получить данные о всех пользователях, 4 - получить данные о пользователе")
	var answer string
	fmt.Scanln(&answer)

	switch answer {
	case "1":
		var login, pass, pass2 string
		db := new(DbWorker)

		fmt.Println("Введите логин")
		fmt.Scanln(&login)
		fmt.Println("Введите пароль")
		fmt.Scanln(&pass)
		fmt.Println("Введите подтверждение пароля")
		fmt.Scanln(&pass2)

		err := db.Add(login, pass, pass2)

		if err != "" {
			fmt.Println(err)
		} else {
			fmt.Println("Успешно")
		}

	case "2":
		db := new(DbWorker)
		fmt.Println("Введите логин для удаления")
		fmt.Scanln(&answer)
		db.DeleteUser(answer)
		fmt.Println("Пользователь удален")
	case "3":
		db := new(DbWorker)
		user_arr := db.SelectUsers()

		for _, value := range user_arr {
			fmt.Println(value[0] + " : " + value[1] + " : " + value[2] + " : " + value[3] + " : " + value[4])
		}
	case "4":
		db := new(DbWorker)
		fmt.Println("Введите логин для получения данных")
		fmt.Scanln(&answer)
		user_arr := db.SelectUser(answer)
		fmt.Println(fmt.Printf("Логин - %s Пароль - %s Подтверждение пароля - %s Результат авторизации - %s Ошибка - %s", user_arr[0], user_arr[1], user_arr[2], user_arr[3], user_arr[4]))
	}

}
