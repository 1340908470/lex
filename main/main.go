package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var INDEX int64 = 0 // token的开始（包括
const FileName = "test.c"

var KEYWORDS = []string{"if", "do", "int", "for",
	"auto", "else", "long", "char", "enum", "void", "case", "goto",
	"short", "float", "union", "const", "while", "break", "double",
	"switch", "struct", "signed", "extern", "static", "sizeof", "return",
	"typedef", "default", "unsigned", "register", "volatile", "continue"}

const (
	// 关键字
	TOEKN_KEYWORD = iota

	// 运算符
	TOKEN_OPERATION

	TOKEN_STRING // 字符串

	// 标识符
	TOKEN_ID // 标识符 变量名、函数名、宏名、结构体名等，由字母、数字、下划线组成，并且首字符必须为字母或下划线

	// 数字
	TOKEN_NUMBER // 无符号数 整数，
)

type Token struct {
	Type int
	Len  int
	Line int

	Desc string // 描述
}

func main() {
	// 读取文件
	file, err := os.Open(FileName)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// 将文件内容整体转为字符串 code
	code := string(data)

	var tokens []Token
	for true {
		if INDEX >= int64(len(code)) {
			break
		}
		if code[INDEX:INDEX+1] == " " || code[INDEX:INDEX+1] == "\r" || code[INDEX:INDEX+1] == "\n" || code[INDEX:INDEX+1] == "\t" {
			INDEX++
			continue
		}
		token := Next(code, INDEX)

		tokens = append(tokens, token)
	}

	for _, token := range tokens {
		fmt.Println(token.Desc)
	}
}

// Next 将code中index指向的字符及其之后的若干个字符转换为Token，index为要读取的第一个字符的索引
func Next(code string, index int64) Token {
	types := -1
	length := 0
	lines := 0
	desc := ""

	/*
		如果判定为注释，则index直接后移至注释结束
	*/
	for len(code[index:]) > 1 && IsAnnotation(code[index:index+2]) {
		anoReg := regexp.MustCompile(`\/\/.*\n`)
		if anoReg == nil {
			fmt.Println("regexp err")
		}
		result := anoReg.FindAllString(code[index:], 1)
		if len(result) != 0 {
			index += int64(len(result[0]))
			INDEX += int64(len(result[0]))
		}
	}

	for code[index:index+1] == " " || code[index:index+1] == "\r" || code[index:index+1] == "\n" || code[index:index+1] == "\t" {
		index++
		INDEX++
	}

	/*
		如果第一个字符为字母，则有可能为：
		- 关键字
		- 标识符
	*/
	if IsAlpha(code[index : index+1]) {
		/* 关键字判别 */
		for _, keyword := range KEYWORDS {
			if strings.Index(code[index:], keyword) == 0 {
				types = TOEKN_KEYWORD
				length = len(keyword)
				desc = "关键字：" + keyword
				break
			}
		}

		/* 标识符判别 */
		// 如果不是关键字，则一定是标识符
		if types == -1 {
			desc = ""
			// 判别标识符的正则表达式：以字母开头或以下划线开头的 包含任意个数的字母、下划线或数字的串
			idReg := regexp.MustCompile(`^([a-z]|[A-Z]|_)\w*`)
			if idReg == nil {
				fmt.Println("regexp err")
			}
			result := idReg.FindAllString(code[index:], 1)
			if len(result) != 0 {
				types = TOKEN_ID
				length = len(result[0])
				desc = "标识符：" + result[0]
			} else {
				panic(nil)
			}
		}

		lines = strings.Count(code[:index], "\n") + 1
	}

	/*
		如果第一个字符为数字，则有可能为：
		- 一定是数字
	*/
	if IsNum(code[index : index+1]) {
		// 匹配规则：数字开始，小数点可选，
		numReg := regexp.MustCompile(`([0-9]+\.?[0-9]*(e|E)-?[0-9]+)|([0-9]*\.?[0-9]+(e|E)[0-9]+)|([0-9]+\.[0-9]*)|([0-9]*\.[0-9]+)|(0(x|X)[0-7]+)|(-?0(b|B)[0-1]+)|([0-9]+)`)
		if numReg == nil {
			fmt.Println("regexp err")
		}
		result := numReg.FindAllString(code[index:], 1)
		if len(result) != 0 {
			types = TOKEN_NUMBER
			length = len(result[0])
			desc = "数字：" + result[0]
		} else {
			panic(nil)
		}
	}

	/*
		如果判定为运算符
	*/
	if (len(code[index:]) > 1 && IsOperation(code[index:index+2])) || IsOperation(code[index:index+1]) {
		OperReg := regexp.MustCompile(`\>=|\<=|!=|\+=|-=|\*=|\/=|%=|&=|\|=|\^=|>>=|<<=|\+\+|\-\-|->|==|=|&&|\+|\-|\*|\\|&|%|<|>|!|\||~|\^|,|\.|\|\||>>|<<|\?|:|\(|\)|\[|\]|\{|\}|;`)
		if OperReg == nil {
			fmt.Println("regexp err")
		}
		result := OperReg.FindAllString(code[index:], 1)
		if len(result) != 0 {
			types = TOKEN_OPERATION
			length = len(result[0])
			desc = "运算符：" + result[0]
		} else {
			panic(nil)
		}
	}

	/*
		如果第一个字符为 " , 则表示为字符串
	*/
	if code[index:index+1] == "\"" {
		codeStr := code[index:]
		// 因为字符串对于引号的匹配是非贪婪的，因此为了在非贪婪模式下略过 \" 的情况，需要先对 \" 进行转换，先将其替换为 ~~，搜索完毕后再换回来
		codeStr = strings.Replace(codeStr, "\\\"", "~~~~~", -1)

		strReg := regexp.MustCompile(`^".*?"`)
		if strReg == nil {
			fmt.Println("regexp err")
		}
		result := strReg.FindAllString(codeStr, 1)

		if len(result) != 0 {
			// 将 result[0] 中的 ~~ 换回 \"
			resultStr := strings.Replace(result[0], "~~~~~", "\\\"", -1)
			types = TOKEN_STRING
			length = len(resultStr)
			desc = "字符串：" + resultStr
		} else {
			panic(nil)
		}
	}

	INDEX += int64(length)
	return Token{
		Type: types,
		Len:  length,
		Line: lines,
		Desc: desc,
	}
}

// IsAlpha 是否为字母
func IsAlpha(s string) bool {
	return (s >= "a" && s <= "z") || (s >= "A" && s <= "Z")
}

// IsNum 是否为数字
func IsNum(s string) bool {
	return s >= "0" && s <= "9"
}

func IsAnnotation(str string) bool {
	if str[0:2] == "/*" || str[0:2] == "//" {
		return true
	} else {
		return false
	}
}

// IsOperation 是否为运算符
func IsOperation(str string) bool {
	s := str[0:1]
	if s == "(" || s == ")" || s == "[" || s == "]" || s == "{" || s == "}" || s == "," ||
		s == "+" || s == "-" || s == "*" || s == "/" || s == ">" || s == "<" || s == "!" ||
		s == "&" || s == "|" || s == "~" || s == "^" || s == "=" || s == "%" || s == ";" {
		// 对于 / 号，需要特殊考虑注释的情况
		if len(str) > 1 && (str[0:2] == "/*" || str[0:2] == "//") {
			return false
		}
		return true
	} else {
		return false
	}
}
