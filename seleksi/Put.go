package seleksi

import (
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mitra
func SeleksiBerkasMahasiswaMagangByMitra(idmahasiswamagang, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.MahasiswaMagang) error {
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(idmahasiswamagang, db)
	if err != nil {
		return err
	}
	_, err = intermoni.GetMagangFromIDByMitra(mahasiswa_magang.Magang.ID, iduser, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.SeleksiBerkas != 0 {
		return SeleksiWewancaraMahasiswaMagangByMitra(idmahasiswamagang, iduser, db, insertedDoc)
	}
	if insertedDoc.SeleksiBerkas != 1 && insertedDoc.SeleksiBerkas != 2 {
		return fmt.Errorf("kesalahan server")
	}
	if insertedDoc.SeleksiBerkas == 2 {
		mahasiswa_magang.SeleksiWewancara = 2
		mahasiswa_magang.Status = 2
	}
	data := bson.M{
		"mahasiswa": bson.M{
			"_id": mahasiswa_magang.Mahasiswa.ID,
		},
		"magang": bson.M{
			"_id": mahasiswa_magang.Magang.ID,
		},
		"pembimbing": bson.M{
			"_id": mahasiswa_magang.Pembimbing.ID,
		},
		"mentor": bson.M{
			"_id": mahasiswa_magang.Mentor.ID,
		},
		"seleksiberkas":    insertedDoc.SeleksiBerkas,
		"seleksiwewancara": mahasiswa_magang.SeleksiWewancara,
		"status": mahasiswa_magang.Status,
	}
	err = intermoni.UpdateOneDoc(idmahasiswamagang, db, "mahasiswa_magang", data)
	if err != nil {
		return err
	}
	mahasiswa, err := intermoni.GetMahasiswaFromID(mahasiswa_magang.Mahasiswa.ID, db)
	if err != nil {
		return err
	}
	user, err := intermoni.GetUserFromID(mahasiswa.Akun.ID, db)
	if err != nil {
		return err
	}
	if (insertedDoc.SeleksiBerkas == 1){
		message := `Selamat ` + mahasiswa.NamaLengkap + `,kamu lolos seleksi berkas\n\nSilahkan lanjut untuk tahap seleksi wewancara,\nAdmin Intern Monitoring`
 		err = intermoni.SendWhatsAppConfirmation(user.Phone, db, message)
		if err != nil {
			return err
		}
		return nil
	}
	if (insertedDoc.SeleksiBerkas == 2){
		message := `Mohon maaf ` + mahasiswa.NamaLengkap + `, kamu tidak lolos seleksi berkas\n\nTetap semangat dan jangan menyerah !!!\nAdmin Intern Monitoring`
 		err = intermoni.SendWhatsAppConfirmation(user.Phone, db, message)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("gagal kirim pesan whastapp")
}

func SeleksiWewancaraMahasiswaMagangByMitra(idmahasiswamagang, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.MahasiswaMagang) error {
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(idmahasiswamagang, db)
	if err != nil {
		return err
	}
	_, err = intermoni.GetMagangFromIDByMitra(mahasiswa_magang.Magang.ID, iduser, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.SeleksiWewancara != 0 {
		return fmt.Errorf("mahasiswa sudah diseleksi")
	}
	if insertedDoc.SeleksiWewancara != 1 && insertedDoc.SeleksiWewancara != 2 {
		return fmt.Errorf("kesalahan server")
	}
	if mahasiswa_magang.SeleksiBerkas != 1 {
		return fmt.Errorf("belum lolos seleksi berkas")
	}
	if insertedDoc.SeleksiWewancara == 2 {
		mahasiswa_magang.Status = 2
	}
	data := bson.M{
		"mahasiswa": bson.M{
			"_id": mahasiswa_magang.Mahasiswa.ID,
		},
		"magang": bson.M{
			"_id": mahasiswa_magang.Magang.ID,
		},
		"pembimbing": bson.M{
			"_id": mahasiswa_magang.Pembimbing.ID,
		},
		"mentor": bson.M{
			"_id": mahasiswa_magang.Mentor.ID,
		},
		"seleksiberkas":    mahasiswa_magang.SeleksiBerkas,
		"seleksiwewancara": insertedDoc.SeleksiWewancara,
		"status": mahasiswa_magang.Status,
	}
	err = intermoni.UpdateOneDoc(idmahasiswamagang, db, "mahasiswa_magang", data)
	if err != nil {
		return err
	}
	mahasiswa, err := intermoni.GetMahasiswaFromID(mahasiswa_magang.Mahasiswa.ID, db)
	if err != nil {
		return err
	}
	user, err := intermoni.GetUserFromID(mahasiswa.Akun.ID, db)
	if err != nil {
		return err
	}
	if (insertedDoc.SeleksiWewancara == 1){
		message := `Selamat ` + mahasiswa.NamaLengkap + `, kamu lolos seleksi wewancara!!!\n\nSilahkan konfirmasi untuk mengambil magang,\nAdmin Intern Monitoring`
 		err = intermoni.SendWhatsAppConfirmation(user.Phone, db, message)
		if err != nil {
			return err
		}
		return nil
	}
	if (insertedDoc.SeleksiWewancara == 2){
		message := `Mohon maaf ` + mahasiswa.NamaLengkap + `, kamu tidak lolos seleksi wewancara\n\nCoba lagi lain kali yaa...\nAdmin Intern Monitoring`
 		err = intermoni.SendWhatsAppConfirmation(user.Phone, db, message)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("gagal kirim pesan whastapp")
}