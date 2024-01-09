package mahasiswa_magang

import (
	"fmt"

	intermoni "github.com/intern-monitoring/backend-intermoni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// by mitra
func TambahMentorMahasiswaMagangByMitra(idmahasiswamagang, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.MahasiswaMagang) error {
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(idmahasiswamagang, db)
	if err != nil {
		return err
	}
	_, err = intermoni.GetMagangFromIDByMitra(mahasiswa_magang.Magang.ID, iduser, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.Status != 1 {
		return fmt.Errorf("mahasiswa belum aktif magang")
	}
	if insertedDoc.Mentor.ID == primitive.NilObjectID {
		return fmt.Errorf("mentor tidak boleh kosong")
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
			"_id": insertedDoc.Mentor.ID,
		},
		"seleksiberkas":    mahasiswa_magang.SeleksiBerkas,
		"seleksiwewancara": mahasiswa_magang.SeleksiWewancara,
		"status": mahasiswa_magang.Status,
	}
	err = intermoni.UpdateOneDoc(idmahasiswamagang, db, "mahasiswa_magang", data)
	if err != nil {
		return err
	}
	return nil
}

// by admin
func TambahPembimbingMahasiswaMagangByAdmin(idmahasiswamagang primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.MahasiswaMagang) error {
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(idmahasiswamagang, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.Status != 1 {
		return fmt.Errorf("mahasiswa belum aktif magang")
	}
	if insertedDoc.Pembimbing.ID == primitive.NilObjectID {
		return fmt.Errorf("pembimbing tidak boleh kosong")
	}
	data := bson.M{
		"mahasiswa": bson.M{
			"_id": mahasiswa_magang.Mahasiswa.ID,
		},
		"magang": bson.M{
			"_id": mahasiswa_magang.Magang.ID,
		},
		"pembimbing": bson.M{
			"_id": insertedDoc.Pembimbing.ID,
		},
		"mentor": bson.M{
			"_id": mahasiswa_magang.Mentor.ID,
		},
		"seleksiberkas":    mahasiswa_magang.SeleksiBerkas,
		"seleksiwewancara": mahasiswa_magang.SeleksiWewancara,
		"status": mahasiswa_magang.Status,
	}
	err = intermoni.UpdateOneDoc(idmahasiswamagang, db, "mahasiswa_magang", data)
	if err != nil {
		return err
	}
	return nil
}

// by mahasiswa
func UpdateStatusMahasiswaMagang(idmahasiswamagang, iduser primitive.ObjectID, db *mongo.Database, insertedDoc intermoni.MahasiswaMagang) error {
	mahasiswa_magang, err := intermoni.GetMahasiswaMagangFromID(idmahasiswamagang, db)
	if err != nil {
		return err
	}
	mahasiswa, err := intermoni.GetMahasiswaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if mahasiswa_magang.Mahasiswa.ID != mahasiswa.ID {
		return fmt.Errorf("kamu bukan pemilik data ini")
	}
	if mahasiswa_magang.Status != 0 {
		return fmt.Errorf("kamu sudah diseleksi")
	}
	if insertedDoc.Status != 1 && insertedDoc.Status != 2  {
		return fmt.Errorf("kesalahan server")
	}
	if insertedDoc.Status == 1 && mahasiswa_magang.SeleksiWewancara != 1{
		return fmt.Errorf("maneh belum lolos seleksi wawancara")
	}
	if insertedDoc.Status == 1 {
		mhs_mgn, err := GetMahasiswaMagangByMahasiswa(iduser, db)
		if err != nil {
			return err
		}
		for _, m := range mhs_mgn {
			m.Status = 2
			if m.ID == mahasiswa_magang.ID {
				m.Status = 1
			}
			data := bson.M{
				"mahasiswa": bson.M{
					"_id": m.Mahasiswa.ID,
				},
				"magang": bson.M{
					"_id": m.Magang.ID,
				},
				"pembimbing": bson.M{
					"_id": m.Pembimbing.ID,
				},
				"mentor": bson.M{
					"_id": m.Mentor.ID,
				},
				"seleksiberkas":    m.SeleksiBerkas,
				"seleksiwewancara": m.SeleksiWewancara,
				"status": m.Status,
			}
			err = intermoni.UpdateOneDoc(m.ID, db, "mahasiswa_magang", data)
			if err != nil {
				return fmt.Errorf("error UpdateStatusMahasiswaMagang update mahasiswa magang tes: %s", err)
			}
		}
		return nil
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
		"seleksiwewancara": mahasiswa_magang.SeleksiWewancara,
		"status": insertedDoc.Status,
	}
	err = intermoni.UpdateOneDoc(idmahasiswamagang, db, "mahasiswa_magang", data)
	if err != nil {
		return err
	}
	return nil
}