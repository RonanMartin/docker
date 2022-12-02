package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Iniciando o servidor REST na porta 8080")
	fmt.Println("Ex.: http://localhost:8080/valida-cpf?numero=91468384066")

	http.HandleFunc("/valida-cpf", ValidaCPF)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func converteCPF(cpfstr string) ([]int, error) {
	var cpfint []int
	for _, numstr := range cpfstr {
		numstr, err := strconv.Atoi(string(numstr))
		if err != nil {
			return nil, err
		}
		cpfint = append(cpfint, numstr)
	}
	return cpfint, nil
}

func recalculaCPF(cpfint []int) []int {
	// duplica os 9 primeiros digitos do cpfint para uma nova variavel
	var cpfpart []int
	cpfpart = append(cpfpart, cpfint[:9]...)

	// calcula digito verificador 1
	multiplicador := 10
	total := 0

	for _, v := range cpfpart {
		total += v * multiplicador
		multiplicador--
	}

	digito1 := 11 - (total % 11)

	switch {
	case digito1 >= 10:
		cpfpart = append(cpfpart, 0)
	case digito1 < 10:
		cpfpart = append(cpfpart, digito1)
	}

	// calcula digito verificador 2
	multiplicador = 11
	total = 0

	for _, v := range cpfpart {
		total += v * multiplicador
		multiplicador--
	}

	digito2 := 11 - (total % 11)

	switch {
	case digito2 >= 10:
		cpfpart = append(cpfpart, 0)
	case digito2 < 10:
		cpfpart = append(cpfpart, digito2)
	}

	return cpfpart
}

func comparaCPFs(cpf1int []int, cpf2int []int) bool {
	cpf1str := make([]string, len(cpf1int))
	for k, v := range cpf1int {
		cpf1str[k] = strconv.Itoa(v)
	}

	cpf2str := make([]string, len(cpf2int))
	for k, v := range cpf2int {
		cpf2str[k] = strconv.Itoa(v)
	}

	cpf1 := strings.Join(cpf1str, "")
	cpf2 := strings.Join(cpf2str, "")

	return cpf1 == cpf2
}

func formataCPF(cpfint []int) string {
	cpfstr := make([]string, len(cpfint))
	for k, v := range cpfint {
		cpfstr[k] = strconv.Itoa(v)
	}

	part1 := strings.Join(cpfstr[:3], "")
	part2 := strings.Join(cpfstr[3:6], "")
	part3 := strings.Join(cpfstr[6:9], "")
	part4 := strings.Join(cpfstr[9:], "")

	return fmt.Sprintf("%s.%s.%s-%s", part1, part2, part3, part4)
}

type Resposta struct {
	Valido    bool   `json:"valido"`
	Formatado string `json:"formatado"`
}

func ValidaCPF(rw http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(rw)

	cpfstr := r.URL.Query().Get("numero")

	cpfint, err := converteCPF(cpfstr)
	if err != nil {
		enc.Encode("CPF invalido")
		return
	}

	recalculado := recalculaCPF(cpfint)

	valido := comparaCPFs(cpfint, recalculado)
	formatado := formataCPF(cpfint)

	rw.Header().Set("Access-Control-Allow-Origin", "*")

	enc.Encode(Resposta{valido, formatado})
}
