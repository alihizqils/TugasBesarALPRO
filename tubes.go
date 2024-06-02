package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const NMAX int = 100 // Konstanta NMAX untuk maks akun
const QMAX int = 100 // Konstanta QMAX untuk maks pertanyaan
const TMAX int = 100 // Konstanta TMAX untuk maks tag

// Tipe bentukan untuk akun
type akun struct {
	username string
	password string
	role     int
}

type acc [NMAX]akun // Tipe alias untuk array of akun dengan ukuran NMAX

// Tipe bentukan pertanyaan
type question struct {
	content  string
	tags     [TMAX]string
	tagCount int
	answer   string
}

type questions [QMAX]question // Tipe alias untuk array of question dengan ukuran QMAX

func main() {
	var A acc
	var Q questions
	var pilihan int
	var i int = 0
	var qIndex int = 0
	var log bool = false
	var currentRole int
	var exit bool = false
	// Perulangan selama pengguna tidak memilih angka 0
	for !exit {

		pilihan = pilihanList(log) // Pemanggilan fungsi untuk menampilkan menu

		switch log {
		// Case bila pengguna belum login
		case false:
			switch pilihan {
			case 1:
				if i < NMAX {
					registrasi(&A, &i) // Pemanggilan prosedur untuk registrasi pengguna
				} else {
					fmt.Println("Maksimum registrasi akun tercapai.")
				}
			case 2:
				if !log {
					login(&A, i, &log, &currentRole) // Pemanggilan prosedur untuk status login dokter dan pasien
				} else {
					fmt.Println("Anda sudah login.")
				}
			case 3:
				if !log {
					guestLogin(&log, &currentRole) // Pemanggilan prosedur untuk status login tamu
				} else {
					fmt.Println("Anda sudah login.")
				}
			case 0:
				exit = true // Pilihan untuk keluar dari program
			default:
				fmt.Println("Pilihan tidak valid, Masukkan angka yang sesuai.")
			}
		// Case bila pengguna sudah login
		case true:
			switch pilihan {
			case 1:
				if log && currentRole == 2 {
					unggahPertanyaan(&Q, &qIndex) // Prosedur untuk pasien mengunggah pertanyaan
				} else {
					fmt.Println("Hanya pasien yang dapat mengunggah pertanyaan, Silakan login terlebih dahulu.")
				}
			case 2:
				if log {
					cariPertanyaan(Q, qIndex) // Fungsi untuk dokter dan pasien mencari pertanyaan berdasarkan tag
				} else {
					fmt.Println("Silakan login terlebih dahulu.")
				}
			case 3:
				if log && currentRole == 1 {
					jawabPertanyaan(&Q, qIndex) // Prosedur untuk dokter menjawab pertanyaan dari pasien
				} else {
					fmt.Println("Hanya dokter yang dapat menjawab pertanyaan, Silakan login terlebih dahulu.")
				}
			case 7:
				if log {
					log = false // status login false untuk logout
					currentRole = 0
					fmt.Println("Berhasil logout.")
				} else {
					fmt.Println("Anda tidak dalam sebuah akun.")
				}
			case 0:
				exit = true // Pilihan untuk keluar dari program
			default:
				fmt.Println("Pilihan tidak valid, Masukkan angka yang sesuai.")
			}
		}
	}
}

func pilihanList(log bool) int {
	// Menampilkan menu utama dan mengembalikan pilihan pengguna
	var pilihan int
	fmt.Println("----------------------------------")
	fmt.Println("|              Menu              |")
	fmt.Println("----------------------------------")
	if !log {
		// Tampilan menu bila pengguna belum login
		fmt.Println("1. Registrasi")
		fmt.Println("2. Login")
		fmt.Println("3. Login sebagai Tamu")
	} else {
		// Tampilan menu bila pengguna sudah login
		fmt.Println("1. Unggah Pertanyaan")
		fmt.Println("2. Cari Pertanyaan")
		fmt.Println("3. Jawab Pertanyaan")
		fmt.Println("7. Logout")
	}
	fmt.Println("0. Quit")
	fmt.Print("Pilih nomor: ")
	fmt.Scan(&pilihan) // Pengguna menginput nomor yang dipilih
	return pilihan
}

func registrasi(A *acc, i *int) {
	// I.S. Array A sebanyak i
	// F.S. Pendaftaran pengguna untuk menjadi dokter atau pasien
	var pilihan int
	fmt.Println("Registrasi sebagai: ")
	fmt.Println("1. Dokter")
	fmt.Println("2. Pasien")
	fmt.Print("Pilih nomor: ")
	fmt.Scan(&pilihan) // Pengguna menginput nomor yang dipilih
	if pilihan == 1 || pilihan == 2 {
		A[*i].role = pilihan
		fmt.Print("Masukkan username: ")
		fmt.Scan(&A[*i].username)
		fmt.Print("Masukkan password: ")
		fmt.Scan(&A[*i].password)
		hasil := cekRegis(*A, *i) // Pengecekan supaya tidak ada username yang sama
		if hasil {
			fmt.Println("Registrasi berhasil.")
			*i++
		} else {
			fmt.Println("Username telah digunakan, Silakan coba kembali di menu registrasi.")
		}
	} else {
		fmt.Println("Pilihan tidak valid, Masukkan angka 1-2.")
	}
}

func cekRegis(A acc, i int) bool {
	// Melakukan binary search pada array yang telah diurutkan berdasarkan username
	low := 0
	high := i - 1

	// Binary search loop
	for low <= high {
		mid := (low + high) / 2
		if A[mid].username == A[i].username {
			return false // Username ditemukan, registrasi gagal
		} else if A[mid].username < A[i].username {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return true // Username tidak ditemukan, registrasi berhasil
}

func login(A *acc, i int, log *bool, currentRole *int) {
	// I.S. terdefinisi array A sebanyak i, status login dan role pengguna
	// F.S. status login antara berhasil dan gagal
	var user, pass string
	var role int
	var coba int = 1
	var gagal int
	var selesai bool = false
	fmt.Println("Login sebagai: ")
	fmt.Println("1. Dokter")
	fmt.Println("2. Pasien")
	fmt.Print("Pilih nomor: ")
	fmt.Scan(&role)
	if role == 1 || role == 2 {
		for coba <= 3 && !selesai {
			fmt.Print("Masukkan username: ")
			fmt.Scan(&user)
			fmt.Print("Masukkan password: ")
			fmt.Scan(&pass)
			*log = cekData(*A, i, role, user, pass) // Fungsi untuk pengecekan apakah akun sudah terdaftar dan benar
			if *log {
				// jika fungsi cekData true maka login berhasil
				*currentRole = role
				fmt.Println("Login berhasil.")
				selesai = true
			} else {
				// jika fungsi cekData false maka login gagal
				fmt.Println("Login gagal.")
				if coba < 3 {
					// Kesempatan 3x untuk pengguna yang masih gagal saat login
					fmt.Println("Silakan coba lagi, Kesempatan anda tersisa", 3-coba)
				} else if coba == 3 {
					fmt.Println("Kesempatan anda habis.")
					fmt.Println("1. Lupa password")
					fmt.Println("0. Kembali")
					fmt.Print("Pilih nomor: ")
					fmt.Scan(&gagal)
					for gagal != 0 && !selesai {
						if gagal == 1 {
							lupaPassword(*A, i) // Pemanggilan fungsi bila pengguna lupa password
							selesai = true
						} else {
							fmt.Println("Pilihan tidak valid, Masukkan angka yang tersedia.")
							fmt.Println("1. Lupa password")
							fmt.Println("0. Kembali")
							fmt.Print("Pilih nomor: ")
							fmt.Scan(&gagal)
						}

					}
				}
				coba++
			}

		}
	} else {
		fmt.Println("Pilihan tidak valid, Masukkan angka 1-2.")
	}
}

func lupaPassword(A acc, i int) {
	// Mengembalikan password dari username yang dicari
	var lupa string
	var selesai bool = false
	fmt.Print("Masukkan username yang lupa password: ")
	fmt.Scan(&lupa)
	// Searching inputan pengguna dengan username dalam array acc
	for !selesai {
		if i == 0 {
			fmt.Println("Tidak ada akun yang terdaftar.")
		} else {
			for j := 0; j < i; j++ {
				if A[j].username == lupa {
					fmt.Println("Password akun anda adalah: ", A[j].password)
					fmt.Println("Silakan coba kembali di menu login.")
				} else {
					fmt.Println("Username tidak ditemukan, Pastikan username dan role anda benar.")
					selesai = true
				}
			}
		}
		selesai = true
	}
}

func guestLogin(log *bool, currentRole *int) {
	// I.S. terdefinisi status login dan role
	// F.S. status login menjadi true supaya tamu bisa masuk ke menu selanjutnya
	*log = true
	*currentRole = 3 // Role 3 represents guest
	fmt.Println("Login sebagai tamu berhasil.")
}

func cekData(A acc, i, role int, user, pass string) bool {
	// Mengembalikan true apabila inputan user dan pass dari pengguna sama dengan salah satu array dari registrasi
	for j := 0; j < i; j++ {
		if A[j].role == role && A[j].username == user && A[j].password == pass {
			return true
		}
	}
	return false
}

func unggahPertanyaan(Q *questions, qIndex *int) {
	// I.S. Terdefinisi array Q dan bilangan qIndex
	// F.S. Mengembalikan pertanyaan dari pengguna

	reader := bufio.NewReader(os.Stdin)
	// Membersihkan buffer untuk memastikan tidak ada input sebelumnya yang tertinggal
	reader.ReadString('\n')

	if *qIndex >= QMAX {
		// Memeriksa apakah jumlah pertanyaan telah mencapai batas maksimum
		fmt.Println("Maksimum jumlah pertanyaan tercapai.")
		return
	}

	fmt.Printf("Masukkan pertanyaan: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		// Menangani kesalahan saat membaca input pertanyaan
		fmt.Println("Error membaca input pertanyaan: ", err)
		return
	}
	// Menyimpan pertanyaan yang telah dibersihkan dari karakter spasi di awal dan akhir
	Q[*qIndex].content = strings.TrimSpace(input)

	fmt.Printf("Masukkan tag: ")
	input, err = reader.ReadString('\n')
	if err != nil {
		// Menangani kesalahan saat membaca input tag
		fmt.Println("Error membaca input tag: ", err)
		return
	}
	// Memisahkan input tag berdasarkan koma dan membersihkan karakter spasi di awal dan akhir
	tags := strings.Split(strings.TrimSpace(input), ",")
	if len(tags) > TMAX {
		// Memeriksa apakah jumlah tag melebihi batas maksimum
		fmt.Println("Jumlah tag melebihi batas maksimum.")
		return
	}

	for i := range tags {
		// Menyimpan setiap tag setelah dibersihkan dari karakter spasi di awal dan akhir
		Q[*qIndex].tags[i] = strings.TrimSpace(tags[i])
	}
	// Menyimpan jumlah tag yang valid
	Q[*qIndex].tagCount = len(tags)

	// Increment qIndex untuk persiapan pertanyaan berikutnya
	*qIndex++
	fmt.Println("Pertanyaan berhasil diunggah.")
}

func cariPertanyaan(Q questions, qIndex int) {
	// I.S. Terdefinisi array Q dan bilangan qIndex
	// F.S. Mencari pertanyaan yang akan dicari dan mengembalikannya

	// Menampilkan semua pertanyaan sebelum diurutkan
	for i := 0; i < qIndex; i++ {
		fmt.Printf("Pertanyaan: %s\n", Q[i].content)
		if Q[i].answer != "" {
			// Menampilkan jawaban jika sudah ada
			fmt.Printf("Jawaban: %s\n", Q[i].answer)
		} else {
			// Menampilkan pesan jika belum ada jawaban
			fmt.Println("Belum ada jawaban.")
		}
	}
	// Membuat jarak agar output lebih mudah dibaca
	fmt.Println()

	// Urutkan pertanyaan yang belum terjawab berdasarkan panjangnya
	sortingPertanyaan(&Q, qIndex)

	reader := bufio.NewReader(os.Stdin)
	// Membersihkan buffer untuk memastikan tidak ada input sebelumnya yang tertinggal
	reader.ReadString('\n')

	fmt.Printf("Masukkan tag untuk mencari pertanyaan: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		// Menangani kesalahan saat membaca input tag
		fmt.Println("Terjadi kesalahan dalam membaca input: ", err)
		return
	}
	// Menghapus spasi di awal dan akhir tag yang dimasukkan
	tag := strings.TrimSpace(input)

	fmt.Println("Hasil pencarian pertanyaan: ")
	found := false

	// Mencari pertanyaan secara berurutan berdasarkan tag
	for i := 0; i < qIndex; i++ {
		tagFound := false
		for j := 0; j < QMAX; j++ {
			// Memeriksa apakah tag cocok dengan salah satu tag di pertanyaan
			if j < len(Q[i].tags) && Q[i].tags[j] == tag {
				tagFound = true
			}
		}
		if tagFound {
			// Menampilkan pertanyaan dan jawaban jika tag ditemukan
			fmt.Printf("Pertanyaan: %s\n", Q[i].content)
			if Q[i].answer != "" {
				fmt.Printf("Jawaban: %s\n", Q[i].answer)
			} else {
				fmt.Println("Belum ada jawaban.")
			}
			found = true
		}
	}

	if !found {
		// Menampilkan pesan jika tidak ada pertanyaan yang ditemukan dengan tag tersebut
		fmt.Println("Tidak ada pertanyaan yang ditemukan dengan tag tersebut.")
	}
}

func sortingPertanyaan(Q *questions, qIndex int) {
	// I.S. Terdefinisi array Q dan bilangan qIndex
	// F.S. Mengurutkan array Q berdasarkan panjang huruf

	//Selection Sort terpanjang ke terpendek
	for i := 0; i < qIndex-1; i++ {
		minIndex := i
		for j := i + 1; j < qIndex; j++ {
			if len(Q[j].content) > len(Q[minIndex].content) {
				minIndex = j
			}
		}
		// Swap the found minimum element with the element at index i
		temp := Q[i]
		Q[i] = Q[minIndex]
		Q[minIndex] = temp
	}
}

func jawabPertanyaan(Q *questions, qIndex int) {
	// I.S. Terdefinisi array Q dan bilangan qIndex
	// F.S. Menampilkan pertanyaan yang belum dijawab dan menyimpan jawaban dari dokter

	if qIndex == 0 {
		// Jika tidak ada pertanyaan dalam array Q
		fmt.Println("Tidak ada pertanyaan.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Daftar pertanyaan yang belum terjawab: ")

	// Urutkan pertanyaan yang belum terjawab berdasarkan panjangnya menggunakan Insertion Sort
	for i := 0; i < qIndex; i++ {
		for j := 0; j < qIndex-1; j++ {
			if len(Q[j].content) > len(Q[j+1].content) {
				// Menukar posisi pertanyaan jika panjangnya lebih besar
				temp := Q[j]
				Q[j] = Q[j+1]
				Q[j+1] = temp
			}
		}
	}

	// Cetak pertanyaan yang belum terjawab
	adaPertanyaanBelumTerjawab := false
	for i := 0; i < qIndex; i++ {
		if Q[i].answer == "" {
			// Menampilkan hanya pertanyaan yang belum memiliki jawaban
			fmt.Printf("%d. %s\n", i+1, Q[i].content)
			adaPertanyaanBelumTerjawab = true
		}
	}

	if !adaPertanyaanBelumTerjawab {
		// Jika semua pertanyaan telah dijawab
		fmt.Println("Tidak ada pertanyaan yang diunggah.")
		return
	}

	fmt.Print("Pilih nomor pertanyaan: ")
	var pertanyaanIndex int
	fmt.Scan(&pertanyaanIndex)
	// Membersihkan buffer setelah membaca input
	reader.ReadString('\n')
	pertanyaanIndex-- // Convert to zero-based index

	// Validasi input nomor pertanyaan
	if pertanyaanIndex < 0 || pertanyaanIndex >= qIndex || Q[pertanyaanIndex].answer != "" {
		fmt.Println("Pilihan tidak valid.")
		return
	}

	fmt.Printf("Masukkan jawaban untuk pertanyaan \"%s\": ", Q[pertanyaanIndex].content)
	answer, err := reader.ReadString('\n')
	if err != nil {
		// Menangani kesalahan saat membaca input jawaban
		fmt.Println("Error membaca input jawaban: ", err)
		return
	}

	// Menyimpan jawaban setelah membersihkan spasi di awal dan akhir
	Q[pertanyaanIndex].answer = strings.TrimSpace(answer)
	fmt.Println("Jawaban berhasil disimpan.")
}
