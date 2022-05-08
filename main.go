package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//Estructura para los comandos
type Comando_General struct {

	//Nombre del comando que ingreso el usuario
	Nombre string

	//Lista de parametros y sus estados
	size       int64
	EstadoSize bool

	fit       string
	EstadoFit bool

	unit       string
	EstadoUnit bool

	path       string
	EstadoPath bool

	typo       string
	EstadoType bool

	name       string
	EstadoName bool

	id       string
	EstadoId bool

	usuario       string
	EstadoUsuario bool

	password       string
	Estadopassword bool

	pwd       string
	Estadopwd bool

	grp       string
	EstadoGrp bool

	EstadoR bool

	cont       string
	EstadoCont bool

	EstadoP bool

	ruta       string
	EstadoRuta bool

	// Lista de comandos para saber cual esta activado
	mkdisk bool
	rmdisk bool
	fdisk  bool
	mount  bool
	mkfs   bool
	login  bool
	logout bool
	mkgrp  bool
	rmgrp  bool
	mkuser bool
	rmusr  bool
	mkfile bool
	mkdir  bool
	pause  bool
	exec   bool
	rep    bool

	//Contador de comandos para que solo venga un comando
	contador_comando int
}

/**************************************************************
	Definicion de structs
***************************************************************/
type Mbr struct {
	Mbr_fecha_cracion [25]byte
	Mbr_tamanio       [12]byte
	Mbr_dsk_signature [12]byte
	Mbr_fit           byte
	Mbr_partition     [4]Particion
}

type Particion struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  [12]byte
	Part_siize  [12]byte
	Part_name   [16]byte
}

type Ebr struct {
	Part_status byte
	Part_fit    byte
	Part_start  [12]byte
	Part_size   [12]byte
	Part_next   [12]byte
	Part_name   [16]byte
}

//variable global
var abcd [1]byte

type listamount struct {
	nombre_partcion string
	id              string
	path_disco      string
	nombre_disco    string
	letra           string
	no_montura      int
	tamanio         int
	inicio          int
}

func ordenarlistamount(entrada []listamount) {
	for x := 0; x < len(entrada); x++ {
		for y := 0; y < len(entrada)-1; y++ {
			actual := entrada[y]
			next := entrada[y+1]
			if actual.id > next.id && actual.no_montura < next.no_montura {
				entrada[y] = next
				entrada[y+1] = actual
			}
		}
	}
}

func imprimirMount(entrada []listamount) {
	fmt.Printf("IMPRIMIENDO LISTA DE MONTURAS \n")
	fmt.Printf("----------------------------\n")
	for x := 0; x < len(entrada); x++ {
		fmt.Printf("ID-MONTURA %s\n", entrada[x].id)
		fmt.Printf("ID-PATH %s \n", entrada[x].path_disco)
		fmt.Printf("ID-PARTICION %s \n", entrada[x].nombre_partcion)
		fmt.Printf("ID-DISCO %s \n", entrada[x].nombre_disco)
		fmt.Printf("----------------------------\n")
	}
}

//lista para montar
var mount []listamount

type listalibre struct {
	size    int
	inicio  int
	fin     int
	tamanio int
	pos     int
}

func ordenarlistalibre(entrada []listalibre) {
	for x := 0; x < len(entrada); x++ {
		for y := 0; y < len(entrada)-1; y++ {
			actual := entrada[y]
			next := entrada[y+1]
			if actual.inicio > next.inicio {
				entrada[y] = next
				entrada[y+1] = actual
			}
		}
	}
}

func ordenarlistalibrepos(entrada []listalibre) {
	for x := 0; x < len(entrada); x++ {
		for y := 0; y < len(entrada)-1; y++ {
			actual := entrada[y]
			next := entrada[y+1]
			if actual.pos > next.pos {
				entrada[y] = next
				entrada[y+1] = actual
			}
		}
	}
}

//lista para reporte
type listarep struct {
	inicio     int
	tamanio    int
	fin        int
	porcentaje int
	tipo       string
}

//INICIALIZAMOS LA ESTRUCTURA GENERAL PARA TENER TODOS LOS COMANDOS
func newComando_General() *Comando_General {
	p := Comando_General{}
	p.Nombre = ""
	p.size = 0
	p.EstadoSize = false
	p.fit = ""
	p.EstadoFit = false
	p.unit = ""
	p.EstadoUnit = false
	p.path = ""
	p.EstadoPath = false
	p.typo = ""
	p.EstadoType = false
	p.name = ""
	p.EstadoName = false
	p.id = ""
	p.EstadoId = false
	p.usuario = ""
	p.EstadoUsuario = false
	p.password = ""
	p.Estadopassword = false
	p.pwd = ""
	p.Estadopwd = false
	p.grp = ""
	p.EstadoGrp = false
	p.EstadoR = false
	p.cont = ""
	p.EstadoCont = false
	p.EstadoP = false
	p.ruta = ""
	p.EstadoRuta = false
	p.mkdisk = false
	p.rmdisk = false
	p.fdisk = false
	p.mount = false
	p.mkfs = false
	p.login = false
	p.logout = false
	p.mkgrp = false
	p.rmgrp = false
	p.mkuser = false
	p.rmusr = false
	p.mkfile = false
	p.mkdir = false
	p.pause = false
	p.exec = false
	p.rep = false
	p.contador_comando = 0
	return &p
}

//verifica si es una letra
func isLetter(caracter byte) int {
	if (caracter >= 'a' && caracter <= 'z') || (caracter >= 'A' && caracter <= 'Z') {
		return 1
	} else {
		return 0
	}
}

func IntToBytes(i int64) []byte {
	if i > 0 {
		return append(big.NewInt(int64(i)).Bytes(), byte(1))
	}
	return append(big.NewInt(int64(i)).Bytes(), byte(0))
}

func BytesToInt(b []byte) int {
	if b[len(b)-1] == 0 {
		return -int(big.NewInt(0).SetBytes(b[:len(b)-1]).Int64())
	}
	return -int(big.NewInt(0).SetBytes(b[:len(b)-1]).Int64())
}

func pausa() {
	fmt.Printf("\n\nSe ha encontrado una pausa \n\n")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nIngrese :\n")
	fmt.Print(">>>>>>>>>:~$ ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	fmt.Printf("\nPausa(): texto ingresado %s\n", text)
}

func Getint(cadena []byte) int {
	var mt string
	for i := 0; i < len(cadena); i++ {
		if cadena[i] != '\x00' {
			mt += string(cadena[i])
		}
	}
	var result int = 0
	n, err := strconv.Atoi(mt)
	if err != nil {
		fmt.Println(mt, "is not an integer.")
		result = -1
	} else {
		result = n
	}
	return result
}

func GetString(cadena []byte) string {
	var mt string
	for i := 0; i < len(cadena); i++ {
		if cadena[i] != '\x00' {
			mt += string(cadena[i])
		}
	}
	return mt
}

/////////////////////////////////////////////////////////////////////////////////
///////CREACION DE ANALIZADOR Y SELECCION DE PARAMETROS
/////////////////////////////////////////////////////////////////////////////////
func Seleccion() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("---------------------")
	fmt.Println("Henrry David Bran Velasquez")
	fmt.Println("201314439")
	fmt.Println("---------------------")
	for {
		fmt.Print("\n\nBienvenido")
		fmt.Print("\nIngrese Comando:")
		fmt.Print(">>>>>>>>>:~$ ")

		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("exit", strings.ToLower(text)) == 0 {
			fmt.Println("Saliendo....")
			break
		} else {
			text += " "
			cadena := []byte(text + "  ")
			Analizador(cadena)
		}
	}
}
func Analizador(cadena []byte) {
	CG := newComando_General()
	var comando string
	var parametro string
	var valorparametro string
	var contador int = 0
	var errorescontador int = 0
	for contador < len(cadena) {
		var entrada byte = cadena[contador]
		if entrada == '#' {
			//COMENTARIO
			fmt.Println("\nAnalizador() : Comentario")
			var comentario string
			for contador < len(cadena) {
				if cadena[contador] != '\n' {
					comentario += string(cadena[contador])
					contador++
				} else {
					break
				}
			}
			fmt.Println("\nComentario: " + comentario)
		} else if entrada == ' ' {
			//ESPACIO VACIO
			contador++
		} else if isLetter(entrada) == 1 {
			//RECONOCEMOS UN COMANDO
			//COMENZAMOS A JUNTAR TODOS LOS CARACTERES QUE SEAN DE TIPO
			//LETRA PARA VER QUE COMANDO ES
			for isLetter(cadena[contador]) == 1 {
				comando += string(cadena[contador])
				contador++
			}
			//YA JUNTOS LOS CARACTERES PROCEDEMOS A VER QUE COMANDO ES
			//POR SI VIENEN MAYUSCULAS Y MINUSCULAS PASAMOS TODO A MINUSCULAS
			comando = strings.ToLower(comando)
			//VERIFICAMOS QUE TIPO DE COMANDO ES Y SI EXISTE
			if comando == "mkdisk" {
				CG.Nombre = comando
				CG.mkdisk = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "rmdisk" {
				CG.Nombre = comando
				CG.rmdisk = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "fdisk" {
				CG.Nombre = comando
				CG.fdisk = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "mount" {
				CG.Nombre = comando
				CG.mount = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "mkfs" {
				CG.Nombre = comando
				CG.mkfs = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "login" {
				CG.Nombre = comando
				CG.login = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "logout" {
				CG.Nombre = comando
				CG.logout = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "mkgrp" {
				CG.Nombre = comando
				CG.mkgrp = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "rmgrp" {
				CG.Nombre = comando
				CG.rmgrp = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "rmusr" {
				CG.Nombre = comando
				CG.rmusr = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "mkfile" {
				CG.Nombre = comando
				CG.mkfile = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "mkdir" {
				CG.Nombre = comando
				CG.mkdir = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "pause" {
				CG.Nombre = comando
				CG.pause = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "exec" {
				CG.Nombre = comando
				CG.exec = true
				CG.contador_comando = CG.contador_comando + 1
			} else if comando == "rep" {
				CG.Nombre = comando
				CG.rep = true
				CG.contador_comando = CG.contador_comando + 1
			} else {
				fmt.Printf("\nAnalizador():ERROR!!! : Comando -> %s \n", comando)
				break
			}
			comando = ""
		} else if entrada == '-' {
			//SUMAMOS AL CONTADOR PARA AGRUPAR LA PALABRA
			contador++
			//concatenamos el parametro
			for cadena[contador] != '=' && cadena[contador] != ' ' {
				parametro += string(cadena[contador])
				contador++
			}
			contador++
			//volvemos a minusculas el parametro
			parametro = strings.ToLower(parametro)
			//comenzamos a tomar los caracteres para el valor del parametro
			//verificamos que el parametro sea de tres tipos uno con unido,otro entre comillas y uno sin valor en el parametro
			if cadena[contador] == '"' {
				//sumamos 1 para saltar las comillas
				contador++
				for cadena[contador] != '"' {
					valorparametro += string(cadena[contador])
					contador++
				}
			} else {
				//verificamos que sea una cadena unida
				for contador < len(cadena) {
					if cadena[contador] != ' ' && cadena[contador] != '\n' && cadena[contador] != '#' {
						valorparametro += string(cadena[contador])
						contador++
					} else {
						break
					}
				}
			}

			//verificamos ya reconocidos los parametros y sus valores
			if parametro == "size" {
				if !CG.EstadoSize {
					intVar, err := strconv.ParseInt(valorparametro, 0, 64)
					if err != nil {
						fmt.Printf("\nAnalizador(): ERRROR !!!!! : valor de parametro %s, no es aceptable -> %s, tipo de error : %s \n", parametro, valorparametro, err)
						return
					} else {
						CG.size = intVar
						CG.EstadoSize = true
					}
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de size")
				}
			} else if parametro == "fit" {
				if !CG.EstadoFit {
					valorparametro = strings.ToUpper(valorparametro)
					if valorparametro == "BF" || valorparametro == "FF" || valorparametro == "WF" {
						CG.fit = valorparametro
						CG.EstadoFit = true
					} else {
						fmt.Printf("\nAnalizador(): ERRROR !!!!! : valor de parametro %s, no es aceptable -> %s\n", parametro, valorparametro)
						errorescontador++
						return
					}
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de fit")
				}
			} else if parametro == "unit" {
				if !CG.EstadoUnit {
					valorparametro = strings.ToUpper(valorparametro)
					if valorparametro == "B" || valorparametro == "K" || valorparametro == "M" {
						CG.unit = valorparametro
						CG.EstadoUnit = true
					} else {
						fmt.Printf("\nAnalizador(): ERRROR !!!!! : valor de parametro %s, no es aceptable -> %s\n", parametro, valorparametro)
						errorescontador++
						return
					}
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de unit")
				}
			} else if parametro == "path" {
				if !CG.EstadoPath {
					CG.path = valorparametro
					CG.EstadoPath = true
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de path")
				}
			} else if parametro == "type" {
				if !CG.EstadoType {
					valorparametro = strings.ToUpper(valorparametro)
					if valorparametro == "P" || valorparametro == "E" || valorparametro == "L" || valorparametro == "FAST" || valorparametro == "FULL" {
						CG.typo = valorparametro
						CG.EstadoType = true
					} else {
						fmt.Printf("\nAnalizador(): ERRROR !!!!! : valor de parametro %s, no es aceptable -> %s\n", parametro, valorparametro)
						errorescontador++
						return
					}
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de type")
				}
			} else if parametro == "name" {
				if !CG.EstadoName {
					CG.name = valorparametro
					CG.EstadoName = true
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de name")
				}
			} else if parametro == "id" {
				if !CG.EstadoId {
					CG.id = valorparametro
					CG.EstadoId = true
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de id")
				}
			} else if parametro == "usuario" {
				if !CG.EstadoUsuario {
					CG.usuario = valorparametro
					CG.EstadoUsuario = true
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de usuario")
				}
			} else if parametro == "password" {
				if !CG.Estadopassword {
					CG.password = valorparametro
					CG.Estadopassword = true
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de password")
				}
			} else if parametro == "pwd" {
				if !CG.Estadopwd {
					CG.pwd = valorparametro
					CG.Estadopwd = true
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de pwd")
				}
			} else if parametro == "grp" {
				if !CG.EstadoGrp {
					CG.grp = valorparametro
					CG.EstadoGrp = true
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de grp")
				}
			} else if parametro == "r" {
				if !CG.EstadoR {
					if valorparametro != "" {
						fmt.Printf("\nAnalizador(): ERRROR !!!!! : valor de parametro %s, no es aceptable -> %s\n", parametro, valorparametro)
						errorescontador++
						return
					} else {
						CG.EstadoR = true
					}
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de R")
				}
			} else if parametro == "cont" {
				if !CG.EstadoCont {
					CG.cont = valorparametro
					CG.EstadoCont = true
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de cont")
				}
			} else if parametro == "p" {
				if !CG.EstadoP {
					if valorparametro != "" {
						fmt.Printf("\nAnalizador(): ERRROR !!!!! : valor de parametro %s, no es aceptable -> %s\n", parametro, valorparametro)
						errorescontador++
						return
					} else {
						CG.EstadoP = true
					}
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de P")
				}
			} else if parametro == "ruta" {
				if !CG.EstadoRuta {
					CG.ruta = valorparametro
					CG.EstadoRuta = true
				} else {
					fmt.Println("\nAnalizador(): ADVERTENCIA !!!!! : parametro duplicado se tomara el primero que se ha ingresado de ruta")
				}
			} else {
				fmt.Printf("\nAnalizador(): ERRROR !!!!! : parametro %s, no es aceptable\n", parametro)
				return
			}
			parametro = ""
			valorparametro = ""
		} else if cadena[contador] == '\n' {
			break
		} else {
			contador++
		}
	}
	//buscamos que parametro es para ejecutar
	if CG.contador_comando == 0 && errorescontador > 0 {
		fmt.Printf("\nAnalizador(): ERRROR !!!!! : existen varios errores no aceptables\n")
		return
	} else {
		SeleccionParametro(*CG)
	}
}
func SeleccionParametro(comando Comando_General) {
	//buscamos cual es el parametro que se ha ingresado
	var parametro = comando.Nombre
	if parametro == "mkdisk" {
		if comando.EstadoSize && comando.EstadoPath {
			CrearDisco(comando)
		} else {
			fmt.Printf("\nSeleccionComando(): ERROR MKDISK: faltan parametros obligatorios")
			return
		}
	} else if parametro == "rmdisk" {
		if comando.EstadoPath {
			EliminarDiscoMsj(comando)
		} else {
			fmt.Printf("\nSeleccionComando(): ERROR RMDISK: faltan parametros obligatorios")
			return
		}
	} else if parametro == "fdisk" {
		if comando.EstadoSize && comando.EstadoPath && comando.EstadoName {
			adminParticion(comando)
		} else {
			fmt.Printf("\nSeleccionComando(): ERROR FDISK: faltan parametros obligatorios")
			return
		}
	} else if parametro == "mount" {
		if comando.EstadoPath && comando.EstadoName {
			MountParticion(comando)
		} else {
			fmt.Printf("\nSeleccionComando(): ERROR MOUNT: faltan parametros obligatorios")
			return
		}
	} else if parametro == "mkfs" {

	} else if parametro == "login" {

	} else if parametro == "logout" {

	} else if parametro == "mkgrp" {

	} else if parametro == "rmgrp" {

	} else if parametro == "mkuser" {

	} else if parametro == "rmusr" {

	} else if parametro == "mkfile" {

	} else if parametro == "mkdir" {

	} else if parametro == "pause" {
		pausa()
	} else if parametro == "exec" {
		if comando.EstadoPath {
			Exec(comando)
		} else {
			fmt.Printf("\nSeleccionComando(): ERROR Exec: faltan parametros obligatorios")
			return
		}
	} else if parametro == "rep" {
		if comando.EstadoName && comando.EstadoPath {
			rep(comando)
		} else {
			fmt.Printf("\nSeleccionComando(): ERROR rep: faltan parametros obligatorios")
			return
		}
	} else {
		fmt.Printf("\nSeleccionComando(): ERROR!!!! : Comando no existe")
		return
	}
}

/////////////////////////////////////////////////////////////////////////////////
///////COMANDO EXEC
/////////////////////////////////////////////////////////////////////////////////
func Exec(comando Comando_General) {
	var aux int = 0
	var valext int

	var path []byte = []byte(comando.path)

	for aux != len(comando.path) {
		if path[aux] == '.' {
			if path[aux+1] == 's' && path[aux+2] == 'c' && path[aux+3] == 'r' && path[aux+4] == 'i' && path[aux+5] == 'p' && path[aux+6] == 't' {
				valext = 1
				break
			}
		}
		aux++
	}

	if valext == 0 {
		fmt.Println("nExec(): ERROR !!!!: extencion del archivo es incorrectar")
	} else {
		file, err := os.Open(comando.path)
		if err != nil {
			fmt.Printf("\nExec(): ERRROR !!!!! : open file %s\n", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		i := 0
		for scanner.Scan() {
			fmt.Printf("\n\nENTRADA\n\n")
			fmt.Printf("\n\nlinea: %d\n", i)
			fmt.Println(scanner.Text())
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\nIngrese [S/N] para analizar linea:\n")
			fmt.Print(">>>>>>>>>:~$ ")
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			if strings.Compare("s", strings.ToLower(text)) == 0 || strings.Compare("si", strings.ToLower(text)) == 0 {
				fmt.Printf("\n-------------------------------INICIO ANALISIS-------------------------------\n\n")
				Analizador([]byte(scanner.Text() + "  "))
				fmt.Printf("\n-------------------------------FINAL ANALISIS--------------------------------\n\n")

			}
			fmt.Println("Saltando linea....")
			i++
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("\nExec(): ERRROR !!!!! : scanner %s\n", err)
			return
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////
//////CREACION Y ELIMINACION DEL DISCO MKDISK Y RMDISK
////////////////////////////////////////////////////////////////////////////////
func CrearDisco(comando Comando_General) {
	var tam int = int(comando.size)
	var unit string = comando.unit
	var ruta string = comando.path

	auxd := []byte(ruta)
	i := 0
	for i < len(auxd) {
		if auxd[i] == '.' && auxd[i+1] == 'd' && auxd[i+2] == 'k' {
			break
		}
		i++
	}
	for auxd[i] != '/' {
		i--
	}
	directorio := string(auxd[0:i])
	//creamos el directorio
	crearDirectorioSiNoExiste(directorio)
	//creamos el archivo
	archivo, err := os.Create(ruta)

	if err != nil {
		fmt.Printf("\nCrearDisco(): ERROR: No se creo el disco no existe ruta \n")
		return
	}

	disco := Mbr{}
	var vacio int64 = 0
	s := &vacio
	var num int = 0

	//tamano del disco

	if tam <= 0 {
		fmt.Printf("\nCrearDisco(): ERROR: tamano del disco es 0 \n")
		return
	}

	if comando.EstadoUnit {
		if strings.Compare(strings.ToUpper(unit), "M") == 0 {
			num = int(tam) * 1024 * 1024
		} else if strings.Compare(strings.ToUpper(unit), "K") == 0 {
			num = int(tam) * 1024
		} else {
			fmt.Printf("\nCrearDisco(): ERROR: Parametro de unit invalido \n")
			return
		}
	} else {
		num = int(tam) * 1024 * 1024
	}

	//Llenando el archivo

	//colocando el primer byte
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	writeNextBytes(archivo, binario.Bytes())

	//situando el cursor en la ultima posicion
	archivo.Seek(int64(num), 0)

	//colocando el ultimo byte para rellenar
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	writeNextBytes(archivo, binario2.Bytes())

	//Regresando el cursor a 0 para escribir el mbr

	//Formando el MBR

	//definimos la fecha del disco
	fechahora := time.Now()
	fechahoraArreglo := strings.Split(fechahora.String(), "")
	fechahoraCadena := ""
	for i := 0; i < 16; i++ {
		fechahoraCadena = fechahoraCadena + fechahoraArreglo[i]
	}
	copy(disco.Mbr_fecha_cracion[:], fechahoraCadena)

	//creamos la variale random
	copy(disco.Mbr_dsk_signature[:], "201314439")

	//ponemos el tipo de ajuste
	if comando.EstadoFit {
		if comando.fit == "BF" {
			disco.Mbr_fit = 'B'
		} else if comando.fit == "FF" {
			disco.Mbr_fit = 'F'
		} else if comando.fit == "WF" {
			disco.Mbr_fit = 'W'
		} else {
			fmt.Printf("\nCrearDisco(): ERROR: fit no valido \n")
			return
		}
	} else {
		disco.Mbr_fit = 'F'
	}

	//cremos la estructura de las particiones
	for i := 0; i < 4; i++ {
		disco.Mbr_partition[i].Part_status = 'L'
		copy(disco.Mbr_partition[i].Part_start[:], "V")
		copy(disco.Mbr_partition[i].Part_siize[:], "V")
		disco.Mbr_partition[i].Part_fit = 'N'
		copy(disco.Mbr_partition[i].Part_name[:], "V")
		disco.Mbr_partition[i].Part_type = 'N'
	}
	defer archivo.Close()
	archivo.Seek(0, 0)

	sp := strconv.Itoa(num)
	copy(disco.Mbr_tamanio[:], sp)

	//Escribiendo el MBR
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, disco)
	writeNextBytes(archivo, binario3.Bytes())
	fmt.Printf("\n********************")
	fmt.Printf("\n %s -> creado con exito", ruta)
	fmt.Printf("\n********************\n")

}
func writeNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}
func EliminarDiscoMsj(comando Comando_General) {
	fmt.Printf("\n\nDESEA ELIMNIAR EL DISCO\n\n")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nIngrese [S/N]:\n")
	fmt.Print(">>>>>>>>>:~$ ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	comparacion := strings.ToLower(text)
	if strings.Compare("s", comparacion) == 0 || strings.Compare("si", comparacion) == 0 {
		EliminarDisco(comando)
	} else {
		fmt.Println("\nEliminarDiscoMsj(): Cancelando operacion de eliminacion")
	}
}
func EliminarDisco(comando Comando_General) {
	fmt.Print(comando.path)
	e := os.Remove(comando.path)
	if e != nil {
		fmt.Println("\nEliminarDisco(): Cancelando operacion de eliminacion surgio un error")
	}
}
func crearDirectorioSiNoExiste(directorio string) {
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.Mkdir(directorio, 0755)
		if err != nil {
			fmt.Println("\nCreardirectorio(): No se pudo crear los directorios")
		}
	}
	fmt.Println("\nCreardirectorio(): Creados exitosamente")
}

/////////////////////////////////////////////////////////////////////////////////
//////CREACION Y MANIPULACION DE LAS PARTICIONES
/////////////////////////////////////////////////////////////////////////////////
func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func adminParticion(comando Comando_General) {
	//variables globales

	primarias := 0
	extendidas := 0
	var iniciopartion int64

	//comenzamos leyendo el disco para sacar info del mbr
	file, err := os.Open(comando.path)
	if err != nil {
		fmt.Printf("\nadminParticion(): Error al leer el disco")
	}
	defer file.Close()

	disco := Mbr{}

	file.Seek(0, 0)
	data := readNextBytes(file, int(unsafe.Sizeof(Mbr{})))
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.LittleEndian, &disco)
	if err != nil {
		fmt.Printf("\nadminParticion(): Error al leer el mbr")
		return
	}
	iniciopartion, err = file.Seek(0, os.SEEK_CUR)
	if err != nil {
		fmt.Printf("\nadminParticion(): Error al leer la posicion")
		return
	}
	aaa := int(unsafe.Sizeof(Mbr{}))
	fmt.Print(aaa)

	// se reconocio MBR
	fmt.Printf("\ntamanio: %s\n", string(disco.Mbr_tamanio[:]))
	fmt.Printf("fit: %c\n", disco.Mbr_fit)
	fmt.Printf("firma: %d\n", disco.Mbr_dsk_signature)
	fmt.Printf("fecha: %s\n", string(disco.Mbr_fecha_cracion[:]))
	//recorremos las particiones para ver las primarias y extendidas
	for i := 0; i < 4; i++ {
		fmt.Println("*********************************")
		fmt.Printf("Fit: %c\n", disco.Mbr_partition[i].Part_fit)
		fmt.Printf("Name: %s\n", string(disco.Mbr_partition[i].Part_name[:]))
		fmt.Printf("size: %s\n", string(disco.Mbr_partition[i].Part_siize[:]))
		fmt.Printf("Start: %s\n", string(disco.Mbr_partition[i].Part_start[:]))
		fmt.Printf("Status: %c\n", disco.Mbr_partition[i].Part_status)
		fmt.Printf("type: %c\n", disco.Mbr_partition[i].Part_type)
		fmt.Println("*********************************")
		nombre := GetString(disco.Mbr_partition[i].Part_name[:])
		if comando.name == nombre {
			fmt.Printf("\nadminParticion(): ERROR: Nombres iguales de particion %s\n", comando.name)
			return
		}
		if disco.Mbr_partition[i].Part_type == 'P' {
			primarias++
		} else if disco.Mbr_partition[i].Part_type == 'E' {
			extendidas++
		}
	}
	fmt.Printf("\nparticiones:\n")
	fmt.Printf("\nprimarias: %d\n", primarias)
	fmt.Printf("\nextendidas: %d\n", extendidas)

	mbrtamanio := Getint(disco.Mbr_tamanio[:])

	//creamos el tamanio de la particion
	var tamanio int64
	if comando.EstadoUnit {
		if comando.unit == "B" {
			tamanio = comando.size
		} else if comando.unit == "K" {
			tamanio = comando.size * 1024
		} else if comando.unit == "M" {
			tamanio = comando.size * 1024 * 1024
		} else {
			fmt.Printf("\nadminParticion(): ERROR: Parametro de unit no valido \n")
			return
		}
	} else {
		tamanio = comando.size * 1024
	}

	//comprobamos el tamanio de la particion > 0
	if tamanio <= 0 {
		fmt.Printf("\nadminParticion(): ERROR: tamanio menor o igual a 0 \n")
		return
	}

	//comprobacion particiones logicas,primarias y extendidas
	if comando.EstadoType {
		if comando.typo == "P" {
			primarias++
		} else if comando.typo == "E" {
			extendidas++
		} else if comando.typo == "L" {
			//logicas
		} else {
			fmt.Printf("\nadminParticion(): ERROR: tipo de particion no es primaria,logica o extendida \n")
			return
		}
	} else {
		primarias++
		comando.typo = "P"
	}

	//comenzamos a llenar las particones
	totalparticiones := primarias + extendidas
	if totalparticiones <= 4 {
		if primarias <= 4 && extendidas <= 1 && (comando.typo == "P" || comando.typo == "E") {
			if int(tamanio) > mbrtamanio {
				fmt.Printf("\nadminParticion(): ERROR: tamanio de la partcion es mas grande que el espacio disponible \n")
				return
			} else {
				var tamanioocupado int = 0
				for i := 0; i < 4; i++ {
					if disco.Mbr_partition[i].Part_status == '1' {
						tamanioocupado += Getint(disco.Mbr_partition[i].Part_siize[:])
					}
				}
				tamaniolibre := mbrtamanio - tamanioocupado - int(iniciopartion)
				if tamaniolibre >= int(tamanio) {
					var lista []listalibre
					var listasizec int = 0
					inicio := int(iniciopartion)
					for i := 0; i < 4; i++ {
						if disco.Mbr_partition[i].Part_status == '1' {
							if Getint(disco.Mbr_partition[i].Part_start[:]) == inicio {
								lista = append(lista, listalibre{listasizec, Getint(disco.Mbr_partition[i].Part_start[:]), 0, Getint(disco.Mbr_partition[i].Part_siize[:]), 0})
								inicio += Getint(disco.Mbr_partition[i].Part_siize[:])
								lista[listasizec].fin = int(inicio)
								listasizec++
							} else {
								lista = append(lista, listalibre{listasizec, Getint(disco.Mbr_partition[i].Part_start[:]), Getint(disco.Mbr_partition[i].Part_start[:]) + Getint(disco.Mbr_partition[i].Part_siize[:]), Getint(disco.Mbr_partition[i].Part_siize[:]), 0})
								inicio = lista[listasizec].fin
								listasizec++
							}
						}
					}
					if listasizec == 0 {
						//creamos la primera particion en el inicio del disco
						if comando.fit == "B" {
							disco.Mbr_partition[0].Part_fit = 'B'
						} else if comando.fit == "F" {
							disco.Mbr_partition[0].Part_fit = 'F'
						} else {
							disco.Mbr_partition[0].Part_fit = 'W'
						}
						//actualizamos el mbr
						copy(disco.Mbr_partition[0].Part_name[:], comando.name)
						copy(disco.Mbr_partition[0].Part_siize[:], strconv.FormatInt(tamanio, 10))
						copy(disco.Mbr_partition[0].Part_start[:], strconv.FormatInt(int64(aaa), 10))
						disco.Mbr_partition[0].Part_status = '1'
						if comando.typo == "P" {
							disco.Mbr_partition[0].Part_type = 'P'
						} else {
							disco.Mbr_partition[0].Part_type = 'E'
						}

						for i := 0; i < 4; i++ {
							fmt.Println("*********************************")
							fmt.Printf("Fit: %c\n", disco.Mbr_partition[i].Part_fit)
							fmt.Printf("Name: %s\n", string(disco.Mbr_partition[i].Part_name[:]))
							fmt.Printf("size: %s\n", string(disco.Mbr_partition[i].Part_siize[:]))
							fmt.Printf("Start: %s\n", string(disco.Mbr_partition[i].Part_start[:]))
							fmt.Printf("Status: %c\n", disco.Mbr_partition[i].Part_status)
							fmt.Printf("type: %c\n", disco.Mbr_partition[i].Part_type)
							fmt.Println("*********************************")
						}

						file2, err := os.OpenFile(comando.path, os.O_RDWR, 0644)
						if err != nil {
							fmt.Printf("\nadminParticion(): Error al leer el disco")
						}
						defer file2.Close()

						file2.Seek(0, 0)
						//Actualizando el MBR
						var binario2 bytes.Buffer
						binary.Write(&binario2, binary.BigEndian, disco)
						writeNextBytes(file2, binario2.Bytes())

						if disco.Mbr_partition[0].Part_type != 'E' {
							fmt.Printf("\nadminParticion(): Se ha creado la particion primaria \n")
							return
						} else {
							ebr := Ebr{}

							copy(ebr.Part_start[:], disco.Mbr_partition[0].Part_start[:])
							ebr.Part_status = '0'
							copy(ebr.Part_next[:], "-1")
							ebr.Part_fit = '0'
							copy(ebr.Part_size[:], "0")
							copy(ebr.Part_name[:], "V")
							file2.Seek(int64(aaa), 0)
							//Actualizando el MBR
							var binario4 bytes.Buffer
							binary.Write(&binario4, binary.BigEndian, ebr)
							writeNextBytes(file2, binario4.Bytes())
							fmt.Printf("\nadminParticion(): Se ha creado la particion extendida y su ebr \n")
							return
						}
					} else {
						ordenarlistalibre(lista)
						inicio = int(iniciopartion)
						estadoinicio := 0
						var listaaux []listalibre
						pos := 0
						auxlista := 0
						auxlista2 := 0
						//
						for inicio < Getint(disco.Mbr_tamanio[:]) {
							if inicio != lista[0].inicio && estadoinicio == 0 {
								//inicio de disco libre
								listaaux = append(listaaux, listalibre{1, inicio, lista[0].inicio, lista[0].inicio - inicio, pos})
								result := lista[0].inicio - inicio
								inicio = inicio + result
								estadoinicio = 1
								pos++
							} else {
								if estadoinicio == 0 {
									inicio = inicio + lista[0].tamanio
								}
								for auxlista < len(lista) {
									auxlista2 = auxlista + 1
									if auxlista2 < len(lista) {
										//har varios en la lista
										juntos := lista[auxlista].fin
										if juntos != lista[auxlista2].inicio {
											listaaux = append(listaaux, listalibre{1, lista[auxlista].fin, lista[auxlista2].inicio, lista[auxlista2].inicio - lista[auxlista].fin, pos})
											result := lista[auxlista2].inicio - lista[auxlista].fin + lista[auxlista].tamanio
											inicio = inicio + result
											pos++
										} else {
											inicio = inicio + lista[auxlista2].tamanio
										}
									} else {
										//este es el ultimo en la lista
										if lista[auxlista].fin != Getint(disco.Mbr_tamanio[:]) {
											listaaux = append(listaaux, listalibre{1, lista[auxlista].fin, Getint(disco.Mbr_tamanio[:]), Getint(disco.Mbr_tamanio[:]) - lista[auxlista].fin, pos})
											result := Getint(disco.Mbr_tamanio[:]) - lista[auxlista].fin + lista[auxlista].tamanio
											inicio = inicio + result
											pos++
										} else {
											inicio = inicio + lista[len(lista)-1].tamanio
										}
									}
									auxlista++
								}
							}
						}

						//terminamos de recorrer el disco
						if len(listaaux) > 0 {
							var estadonodo int = -5
							if disco.Mbr_fit == 'B' {
								ordenarlistalibre(listaaux)
								for i := 0; i < len(listaaux); i++ {
									if listaaux[i].tamanio >= int(tamanio) {
										estadonodo = i
										break
									}
								}
							} else if disco.Mbr_fit == 'W' {
								ordenarlistalibre(listaaux)
								for i := len(listaaux); i >= 0; i-- {
									if listaaux[i].tamanio >= int(tamanio) {
										estadonodo = i
										break
									}
								}
							} else {
								ordenarlistalibrepos(listaaux)
								for i := 0; i < len(listaaux); i++ {
									if listaaux[i].tamanio >= int(tamanio) {
										estadonodo = i
										break
									}
								}
							}

							if estadonodo != -5 {
								for i := 0; i < 4; i++ {
									if disco.Mbr_partition[i].Part_status == 'L' {
										if comando.EstadoFit {
											if comando.fit == "B" {
												disco.Mbr_partition[i].Part_fit = 'B'
											} else if comando.fit == "F" {
												disco.Mbr_partition[i].Part_fit = 'F'
											} else {
												disco.Mbr_partition[i].Part_fit = 'W'
											}
										} else {
											disco.Mbr_partition[i].Part_fit = 'W'
										}

										if comando.typo == "P" {
											disco.Mbr_partition[i].Part_type = 'P'
										} else {
											disco.Mbr_partition[i].Part_type = 'E'
										}

										//actualizamos el mbr
										copy(disco.Mbr_partition[i].Part_name[:], comando.name)
										siize := strconv.FormatInt(tamanio, 10)
										copy(disco.Mbr_partition[i].Part_siize[:], siize)
										partstar := strconv.Itoa(listaaux[estadonodo].inicio)
										disco.Mbr_partition[i].Part_status = '1'
										copy(disco.Mbr_partition[i].Part_start[:], partstar)
										if comando.typo == "P" {
											disco.Mbr_partition[i].Part_type = 'P'
										} else {
											disco.Mbr_partition[i].Part_type = 'E'
										}

										file3, err := os.OpenFile(comando.path, os.O_RDWR, 0644)
										if err != nil {
											fmt.Printf("\nadminParticion(): Error al leer el disco")
										}
										defer file3.Close()

										file3.Seek(0, 0)
										//Actualizando el MBR
										var binario3 bytes.Buffer
										binary.Write(&binario3, binary.BigEndian, disco)
										writeNextBytes(file3, binario3.Bytes())

										if disco.Mbr_partition[i].Part_type != 'E' {
											fmt.Printf("\nadminParticion(): Se ha creado la particion primaria \n")
											return
										} else {
											ebr := Ebr{}

											copy(ebr.Part_start[:], disco.Mbr_partition[i].Part_start[:])
											ebr.Part_status = '0'
											copy(ebr.Part_next[:], "-1")
											ebr.Part_fit = '0'
											copy(ebr.Part_size[:], "0")
											copy(ebr.Part_name[:], "V")
											file3.Seek(int64(Getint(disco.Mbr_partition[i].Part_start[:])), 0)
											//Actualizando el MBR
											var binario4 bytes.Buffer
											binary.Write(&binario4, binary.BigEndian, ebr)
											writeNextBytes(file3, binario4.Bytes())
											fmt.Printf("\nadminParticion(): Se ha creado la particion extendida y su ebr \n")
											return
										}
									}
								}
							}
						}

					}
				}
			}
		} else if extendidas == 1 && comando.typo == "L" {
			var k int
			for k = 0; k < 4; k++ {
				if disco.Mbr_partition[k].Part_type == 'E' {
					break
				}
			}

			if int(tamanio) > Getint(disco.Mbr_partition[k].Part_siize[:]) {
				fmt.Printf("\nadminParticion(): ERROR: el tamanio de la partcion excede el espacio disponible \n")
				return
			} else {

				var lista []listalibre
				ebr := Ebr{}
				ebr.Part_status = '1'
				file.Seek(int64(Getint(disco.Mbr_partition[k].Part_start[:])), 0)
				data := readNextBytes(file, int(unsafe.Sizeof(Ebr{})))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.LittleEndian, &ebr)
				if err != nil {
					fmt.Printf("\nadminParticion(): Error al leer el ebr")
					return
				}
				if ebr.Part_status != '0' {
					var nodo listalibre
					nodo.inicio = Getint(ebr.Part_start[:])
					nodo.fin = Getint(ebr.Part_start[:]) + Getint(ebr.Part_size[:])
					nodo.tamanio = Getint(ebr.Part_size[:])
					nodo.pos = 0
					if GetString(ebr.Part_name[:]) == comando.name {
						fmt.Printf("\nadminParticion(): Error : nombre repetido en el ebr")
						return
					}
					lista = append(lista, nodo)
				}
				mkll := GetString(ebr.Part_next[:])
				if mkll != "-1" {
					kl := Getint(ebr.Part_next[:])
					file.Seek(int64(kl), 0)
					for ebr.Part_status != '0' {
						data := readNextBytes(file, int(unsafe.Sizeof(Ebr{})))
						buffer := bytes.NewBuffer(data)
						err = binary.Read(buffer, binary.LittleEndian, &ebr)
						if err != nil {
							fmt.Printf("\nadminParticion(): Error al leer el ebr")
							return
						}
						var nodo listalibre
						nodo.inicio = Getint(ebr.Part_start[:])
						nodo.fin = Getint(ebr.Part_start[:]) + Getint(ebr.Part_size[:])
						nodo.tamanio = Getint(ebr.Part_size[:])
						nodo.pos = 0
						if ebr.Part_status != '0' {
							if GetString(ebr.Part_name[:]) == comando.name {
								fmt.Printf("\nadminParticion(): Error : nombre repetido en el ebr")
								return
							}
							lista = append(lista, nodo)
						}
						mkl := GetString(ebr.Part_next[:])
						if mkl != "-1" {
							kl := Getint(ebr.Part_next[:])
							file.Seek(int64(kl), 0)
						} else {
							break
						}
					}
				}

				final := len(lista)
				var espacioocupado int
				//contamos el espacio disponible
				for i := 0; i < len(lista); i++ {
					espacioocupado += lista[i].tamanio
				}
				espaciolibreebr := Getint(disco.Mbr_partition[k].Part_siize[:]) - espacioocupado

				if tamanio <= int64(espaciolibreebr) {
					file3, err := os.OpenFile(comando.path, os.O_RDWR, 0644)
					if err != nil {
						fmt.Printf("\nadminParticion(): Error al leer el disco")
					}
					defer file3.Close()
					if final == 0 {
						file3.Seek(int64(aaa), 0)
						copy(ebr.Part_start[:], disco.Mbr_partition[k].Part_start[:])
						ebr.Part_status = '1'
						copy(ebr.Part_next[:], "-1")
						if comando.EstadoFit {
							if comando.fit == "F" {
								ebr.Part_fit = 'F'
							} else if comando.fit == "B" {
								ebr.Part_fit = 'B'
							} else {
								ebr.Part_fit = 'W'
							}
						} else {
							ebr.Part_fit = 'W'
						}
						copy(ebr.Part_size[:], strconv.FormatInt(tamanio, 10))
						copy(ebr.Part_name[:], []byte(comando.name))
						//Actualizando el EBR
						var binario4 bytes.Buffer
						binary.Write(&binario4, binary.BigEndian, ebr)
						writeNextBytes(file3, binario4.Bytes())

						fmt.Printf("\nadminParticion(): Se ha creado el ebr \n")
						return
					} else {
						ebrnuevo := Ebr{}
						inisal := lista[final-1].inicio + lista[final-1].tamanio
						ini := strconv.FormatInt(int64(inisal), 10)
						copy(ebrnuevo.Part_start[:], ini)
						ebrnuevo.Part_status = '1'
						copy(ebrnuevo.Part_next[:], "-1")
						if comando.EstadoFit {
							if comando.fit == "F" {
								ebrnuevo.Part_fit = 'F'
							} else if comando.fit == "B" {
								ebrnuevo.Part_fit = 'B'
							} else {
								ebrnuevo.Part_fit = 'W'
							}
						} else {
							ebrnuevo.Part_fit = 'W'
						}
						copy(ebrnuevo.Part_size[:], strconv.FormatInt(tamanio, 10))
						copy(ebrnuevo.Part_name[:], []byte(comando.name))

						//Actualizando el EBR
						file3.Seek(int64(inisal), 0)
						var binario4 bytes.Buffer
						binary.Write(&binario4, binary.BigEndian, ebrnuevo)
						writeNextBytes(file3, binario4.Bytes())

						fmt.Printf("\nadminParticion(): Se ha creado el ebr \n")

						//actualizamos el ebr anterior

						copy(ebr.Part_next[:], ebrnuevo.Part_start[:])
						//Actualizando el EBR
						Getint(ebr.Part_start[:])
						file3.Seek(int64(Getint(ebr.Part_start[:])), 0)
						var binario5 bytes.Buffer
						binary.Write(&binario5, binary.BigEndian, ebr)
						writeNextBytes(file3, binario5.Bytes())

						fmt.Printf("\nadminParticion(): se ha actualizado el ebr anterior \n")
						return

					}
				} else {
					fmt.Printf("\nadminParticion(): ERROR: ya no hay mucho espacio libre \n")
					return
				}
			}
		} else {
			fmt.Printf("\nadminParticion(): ERROR: ya existen 2 particiones extendidas y/o no se acepta este tipo de particion \n")
			return
		}
	} else {
		fmt.Printf("\nadminParticion(): ERROR: ya existen 4 particiones principales \n")
		return
	}

}

/////////////////////////////////////////////////////////////////////////////////
//////CREACION Y MANIPULACION DE MONTAJE
/////////////////////////////////////////////////////////////////////////////////
func MountParticion(comando Comando_General) {
	var tamaniopartcion int
	var inicioparticion int

	//abrimos el archivo
	file, err := os.Open(comando.path)
	if err != nil {
		fmt.Printf("\nMountParticion(): Error al leer el disco")
	}
	defer file.Close()

	disco := Mbr{}

	file.Seek(0, 0)
	data := readNextBytes(file, int(unsafe.Sizeof(Mbr{})))
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.LittleEndian, &disco)
	if err != nil {
		fmt.Printf("\nMountParticion(): Error al leer el mbr")
		return
	}

	var encontrado int
	var pos int = -1

	for i := 0; i < 4; i++ {
		if GetString(disco.Mbr_partition[i].Part_name[:]) == comando.name {
			if disco.Mbr_partition[i].Part_type == 'E' {
				fmt.Printf("\nMountParticion(): Error: no se puede montar una partcion extendida")
				return
			} else {
				encontrado++
				tamaniopartcion = Getint(disco.Mbr_partition[i].Part_siize[:])
				inicioparticion = Getint(disco.Mbr_partition[i].Part_start[:])
				break
			}
		}
	}

	if encontrado == 0 {
		for i := 0; i < 4; i++ {
			if disco.Mbr_partition[i].Part_type == 'E' {
				pos = i
				break
			}
		}
	}

	if encontrado == 0 && pos >= 0 {

		ebr := Ebr{}
		ebr.Part_status = '1'
		file.Seek(int64(Getint(disco.Mbr_partition[pos].Part_start[:])), 0)
		data := readNextBytes(file, int(unsafe.Sizeof(Ebr{})))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.LittleEndian, &ebr)
		if err != nil {
			fmt.Printf("\nMountParticion(): Error al leer el ebr")
			return
		}

		if ebr.Part_status != '0' {
			if GetString(ebr.Part_name[:]) == comando.name {
				encontrado++
				tamaniopartcion = Getint(ebr.Part_size[:]) - int(unsafe.Sizeof(Ebr{}))
				inicioparticion = int(unsafe.Sizeof(Ebr{})) + Getint(ebr.Part_start[:])
			}
		}
		mkll := GetString(ebr.Part_next[:])
		if mkll != "-1" {
			kl := Getint(ebr.Part_next[:])
			file.Seek(int64(kl), 0)
			for ebr.Part_status != '0' {
				data := readNextBytes(file, int(unsafe.Sizeof(Ebr{})))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.LittleEndian, &ebr)
				if err != nil {
					fmt.Printf("\nMountParticion(): Error al leer el ebr")
					return
				}
				if ebr.Part_status != '0' {
					if GetString(ebr.Part_name[:]) == comando.name {
						encontrado++
						tamaniopartcion = Getint(ebr.Part_size[:]) - int(unsafe.Sizeof(Ebr{}))
						inicioparticion = int(unsafe.Sizeof(Ebr{})) + Getint(ebr.Part_start[:])
					}
				}
				mkl := GetString(ebr.Part_next[:])
				if mkl != "-1" {
					kl := Getint(ebr.Part_next[:])
					file.Seek(int64(kl), 0)
				} else {
					break
				}
			}
		}
	}

	if encontrado == 1 {
		if len(mount) == 0 {
			var nodomount listamount
			var contador int = 1
			var id bytes.Buffer
			id.WriteString("39")
			id.WriteString(strconv.Itoa(contador))
			idaux := string(abcd[0])
			id.WriteString(idaux)
			nodomount.id = id.String()
			nodomount.letra = string(abcd[0])
			nodomount.no_montura = contador

			pathh := []byte(comando.path)
			var indice int = -1
			for i := 0; i < len(pathh); i++ {
				if pathh[i] == '.' {
					indice = i - 1
					break
				}
			}

			var inicion int
			for j := indice; j >= 0; j-- {
				if pathh[j] == '/' {
					inicion = j + 1
					break
				}
			}

			var nameparti string
			for i := 0; i <= indice; i++ {
				if i >= inicion && i <= indice {
					nameparti += string(pathh[i])
				}
			}

			nodomount.inicio = inicioparticion
			nodomount.tamanio = tamaniopartcion
			nodomount.nombre_disco = nameparti
			nodomount.nombre_partcion = comando.name
			nodomount.path_disco = comando.path
			mount = append(mount, nodomount)
			ordenarlistamount(mount)
			imprimirMount(mount)
		} else {
			for i := 0; i < len(mount); i++ {
				if mount[i].nombre_partcion == comando.name && comando.path == mount[i].path_disco {
					fmt.Printf("\nMountParticion(): ya esta montada la particion")
					return
				}
			}

			pathh := []byte(comando.path)
			var indice int = -1
			for i := 0; i < len(pathh); i++ {
				if pathh[i] == '.' {
					indice = i - 1
					break
				}
			}

			var inicion int
			for j := indice; j >= 0; j-- {
				if pathh[j] == '/' {
					inicion = j + 1
					break
				}
			}

			var nameparti string
			for i := 0; i <= indice; i++ {
				if i >= inicion && i <= indice {
					nameparti += string(pathh[i])
				}
			}
			var pos int = 0
			var aux3 int
			for i := 0; i < len(mount); i++ {
				if mount[i].nombre_disco == nameparti {
					pos++
					aux3 = i
				}
			}

			if pos == 0 && pos < 9 {

				var nodomount listamount
				var contador int = 1
				var id bytes.Buffer
				id.WriteString("39")
				abcd[0] = abcd[0] + '\x01'
				idaux := strconv.Itoa(contador)
				id.WriteString(idaux)
				id.WriteString(string(abcd[0]))
				nodomount.id = id.String()
				nodomount.letra = string(abcd[0])
				nodomount.no_montura = contador

				pathh := []byte(comando.path)
				var indice int = -1
				for i := 0; i < len(pathh); i++ {
					if pathh[i] == '.' {
						indice = i - 1
						break
					}
				}

				var inicion int
				for j := indice; j >= 0; j-- {
					if pathh[j] == '/' {
						inicion = j + 1
						break
					}
				}

				var nameparti string
				for i := 0; i <= indice; i++ {
					if i >= inicion && i <= indice {
						nameparti += string(pathh[i])
					}
				}

				nodomount.inicio = inicioparticion
				nodomount.tamanio = tamaniopartcion
				nodomount.nombre_disco = nameparti
				nodomount.nombre_partcion = comando.name
				nodomount.path_disco = comando.path
				mount = append(mount, nodomount)
				ordenarlistamount(mount)
				imprimirMount(mount)
			} else {
				var nodomount listamount
				contador := mount[aux3].no_montura + 1
				var id bytes.Buffer
				id.WriteString("39")
				id.WriteString(strconv.Itoa(contador))
				idaux := mount[aux3].letra
				id.WriteString(idaux)
				nodomount.id = id.String()
				nodomount.letra = string(idaux)
				nodomount.no_montura = contador

				pathh := []byte(comando.path)
				var indice int = -1
				for i := 0; i < len(pathh); i++ {
					if pathh[i] == '.' {
						indice = i - 1
						break
					}
				}

				var inicion int
				for j := indice; j >= 0; j-- {
					if pathh[j] == '/' {
						inicion = j + 1
						break
					}
				}

				var nameparti string
				for i := 0; i <= indice; i++ {
					if i >= inicion && i <= indice {
						nameparti += string(pathh[i])
					}
				}

				nodomount.inicio = inicioparticion
				nodomount.tamanio = tamaniopartcion
				nodomount.nombre_disco = nameparti
				nodomount.nombre_partcion = comando.name
				nodomount.path_disco = comando.path
				mount = append(mount, nodomount)
				ordenarlistamount(mount)
				imprimirMount(mount)
			}
		}
	} else {
		fmt.Printf("\nMountParticion(): No existe la partcion en el sistema")
		return
	}

}

/////////////////////////////////////////////////////////////////////////////////
//////CREACION REP
/////////////////////////////////////////////////////////////////////////////////
func repdisk(nodomount int, comando Comando_General) {
	//creamos el directorio
	var ruta string = comando.path

	auxd := []byte(ruta)
	i := 0
	for i < len(auxd) {
		if auxd[i] == '.' {
			break
		}
		i++
	}

	for auxd[i] != '/' {
		i--
	}
	directorio := string(auxd[0:i])
	//creamos el directorio
	crearDirectorioSiNoExiste(directorio)

	//abrimos el archivo

	file, err := os.Open(mount[nodomount].path_disco)
	if err != nil {
		fmt.Printf("\nrep(): ERROR: no se abrio el archivo de ruta \n")
		return
	}
	defer file.Close()

	disco := Mbr{}

	file.Seek(0, 0)
	data := readNextBytes(file, int(unsafe.Sizeof(Mbr{})))
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.LittleEndian, &disco)
	if err != nil {
		fmt.Printf("\nadminParticion(): Error al leer el mbr")
		return
	}
	var lista []listarep
	tamaniototal := Getint(disco.Mbr_tamanio[:])
	var extendida int = -1
	var primarias int = -1
	for i := 0; i < 4; i++ {
		if disco.Mbr_partition[i].Part_status == '1' && disco.Mbr_partition[i].Part_type == 'P' {
			var nodo listarep
			nodo.inicio = Getint(disco.Mbr_partition[i].Part_start[:])
			nodo.fin = Getint(disco.Mbr_partition[i].Part_start[:]) + Getint(disco.Mbr_partition[i].Part_siize[:])
			nodo.tamanio = Getint(disco.Mbr_partition[i].Part_siize[:])
			nodo.porcentaje = (Getint(disco.Mbr_partition[i].Part_siize[:]) * 100) / tamaniototal
			nodo.tipo = string(disco.Mbr_partition[i].Part_type)
			lista = append(lista, nodo)
			primarias = i
		}

		if disco.Mbr_partition[i].Part_type == 'E' {
			extendida = i
		}
	}

	if extendida != -1 {
		ebr := Ebr{}
		ebr.Part_status = '1'
		file.Seek(int64(Getint(disco.Mbr_partition[extendida].Part_start[:])), 0)
		data := readNextBytes(file, int(unsafe.Sizeof(Ebr{})))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.LittleEndian, &ebr)
		if err != nil {
			fmt.Printf("\nrep(): Error al leer el ebr")
			return
		}
		if ebr.Part_status != '0' {
			var nodo listarep
			nodo.inicio = Getint(ebr.Part_start[:])
			nodo.fin = Getint(ebr.Part_start[:]) + Getint(ebr.Part_size[:])
			nodo.tamanio = Getint(ebr.Part_size[:])
			nodo.porcentaje = (Getint(ebr.Part_size[:]) * 100) / tamaniototal
			nodo.tipo = "L"
			lista = append(lista, nodo)
		}
		mkll := GetString(ebr.Part_next[:])
		if mkll != "-1" {
			kl := Getint(ebr.Part_next[:])
			file.Seek(int64(kl), 0)
			for ebr.Part_status != '0' {
				data := readNextBytes(file, int(unsafe.Sizeof(Ebr{})))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.LittleEndian, &ebr)
				if err != nil {
					fmt.Printf("\nadminParticion(): Error al leer el ebr")
					return
				}
				var nodo listarep
				nodo.inicio = Getint(ebr.Part_start[:])
				nodo.fin = Getint(ebr.Part_start[:]) + Getint(ebr.Part_size[:])
				nodo.tamanio = Getint(ebr.Part_size[:])
				nodo.porcentaje = (Getint(ebr.Part_size[:]) * 100) / tamaniototal
				nodo.tipo = "L"
				if ebr.Part_status != '0' {
					if GetString(ebr.Part_name[:]) == comando.name {
						fmt.Printf("\nadminParticion(): Error : nombre repetido en el ebr")
						return
					}
					lista = append(lista, nodo)
				}
				mkl := GetString(ebr.Part_next[:])
				if mkl != "-1" {
					kl := Getint(ebr.Part_next[:])
					file.Seek(int64(kl), 0)
				} else {
					break
				}
			}
		}
		extenidadatotal := Getint(ebr.Part_start[:]) + Getint(ebr.Part_size[:])
		if extenidadatotal < (Getint(disco.Mbr_partition[extendida].Part_start[:]) + Getint(disco.Mbr_partition[extendida].Part_siize[:])) {
			var nodo listarep
			nodo.inicio = Getint(ebr.Part_start[:]) + Getint(ebr.Part_size[:])
			nodo.fin = (Getint(disco.Mbr_partition[extendida].Part_start[:]) + Getint(disco.Mbr_partition[extendida].Part_siize[:]))
			nodo.tamanio = (Getint(disco.Mbr_partition[extendida].Part_start[:]) + Getint(disco.Mbr_partition[extendida].Part_siize[:])) - Getint(ebr.Part_start[:]) + Getint(ebr.Part_size[:])
			nodo.porcentaje = (nodo.tamanio * 100) / tamaniototal
			nodo.tipo = "EF"
			lista = append(lista, nodo)
		}
	} else {
		var nodo listarep
		nodo.inicio = Getint(disco.Mbr_partition[primarias].Part_start[:]) + Getint(disco.Mbr_partition[primarias].Part_siize[:])
		nodo.fin = tamaniototal
		nodo.tamanio = tamaniototal - nodo.inicio
		nodo.porcentaje = (nodo.tamanio * 100) / tamaniototal
		nodo.tipo = "PF"
		lista = append(lista, nodo)
	}

	var k int = 0
	for k < len(auxd) {
		if auxd[k] == '.' {

			break
		}
		k++
	}
	nruta := string(auxd[0:k]) + ".dot"
	nimage := string(auxd[0:k]) + ".png"

	var contenido string
	contenido += "digraph {\n"
	contenido += "\ttbl [\n"
	contenido += "\t\tshape=plaintext\n"
	contenido += "\t\tlabel=<\n"
	contenido += "\t\t<table border='0' cellborder='5' color='blue' cellspacing='0'>\n"
	contenido += "\t\t\t<tr>\n"
	contenido += "\t\t\t\t<td>MBR</td>\n"
	for aux3 := 0; aux3 < len(lista); aux3++ {
		if lista[aux3].tipo == "P" {
			contenido += "\t\t\t\t<td rowspan='2'>Primaria<br/>" + strconv.Itoa(lista[aux3].porcentaje) + "%" + " del disco </td>\n"
		} else if lista[aux3].tipo == "L" {
			contenido += "\t\t\t\t<td cellpadding='1'>\n"
			contenido += "\t\t\t\t\t<table color='orange' cellspacing='0'>\n"
			contenido += "\t\t\t\t\t<tr><td colspan=\"20\">Extendida</td></tr>\n"
			contenido += "\t\t\t\t\t<tr>\n"
			for aux3 < len(lista) {
				if lista[aux3].tipo == "L" {
					contenido += "\t\t\t\t\t<td>EBR</td>\n"
					contenido += "\t\t\t\t\t<td rowspan='2'>Logica<br/>" + strconv.Itoa(lista[aux3].porcentaje) + "%" + "del disco </td>\n"
				} else if lista[aux3].tipo == "EL" {
					contenido += "\t\t\t\t\t<td rowspan='2'>Libre<br/>" + strconv.Itoa(lista[aux3].porcentaje) + "%" + "del disco </td>\n"
				} else if lista[aux3].tipo == "P" {
					break
				}
				aux3++
			}
			contenido += "\t\t\t\t\t</tr>\n"
			contenido += "\t\t\t\t\t</table>\n"
			contenido += "\t\t\t\t</td>\n"
		} else if lista[aux3].tipo == "EP" {
			contenido += "\t\t\t\t<td rowspan='2'>Libre<br/>" + strconv.Itoa(lista[aux3].porcentaje) + "%" + "del disco </td>\n"
		}
	}
	contenido += "\t\t</tr>\n"
	contenido += "\t\t</table>\n"
	contenido += "\t>];\n"
	contenido += "}"

	b := []byte(contenido)
	err = ioutil.WriteFile(nruta, b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-Tpng", nruta).Output()
	mode := int(0777)
	ioutil.WriteFile(nimage, cmd, os.FileMode(mode))

	fmt.Printf("\nrep(): se ha creado el reporte exitosamente\n")

}

func rep(comando Comando_General) {
	var aux2 = -1
	for i := 0; i < len(mount); i++ {
		if mount[i].id == comando.id {
			aux2 = i
			break
		}
	}
	if aux2 != -1 {
		if comando.name == "disk" {
			repdisk(aux2, comando)
		} else {
			fmt.Printf("\nrep(): falta por hacer\n")
			return
		}
	} else {
		fmt.Printf("\nrep(): Error no exite esta id de particion montada\n")
		return
	}
}

func mks(nodomount listamount, comando Comando_General) {

}

/////////////////////////////////////////////////////////////////////////////////
//////MAIN
/////////////////////////////////////////////////////////////////////////////////
func main() {
	abcd[0] = 'a'
	fmt.Println("Ingresando al Sistema ....")
	Seleccion()
}
