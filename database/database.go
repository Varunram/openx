package database

import (
	"log"
	"os"

	consts "github.com/OpenFinancing/openfinancing/consts"
	utils "github.com/OpenFinancing/openfinancing/utils"
	"github.com/boltdb/bolt"
)

// package database contains base level stuff that will be required by all the
// sub platforms in the platform of platforms model. Currently contains Investors,
// Recipients and Users, but can be expanded to contain even stages, if that's deemed
// to be common across platforms
// define the name of the buckets that we interact with.
var ProjectsBucket = []byte("Projects")
var InvestorBucket = []byte("Investors")
var RecipientBucket = []byte("Recipients")
var ContractorBucket = []byte("Contractors")
var UserBucket = []byte("Users")
var BondBucket = []byte("Bonds")
var CoopBucket = []byte("Coop")
var InspectorBucket = []byte("Inspector")

func CreateHomeDir() {
	if _, err := os.Stat(consts.HomeDir); os.IsNotExist(err) {
		// directory does not exist, create one
		log.Println("Creating home directory")
		os.MkdirAll(consts.HomeDir, os.ModePerm)
	}
	if _, err := os.Stat(consts.DbDir); os.IsNotExist(err) {
		os.MkdirAll(consts.DbDir, os.ModePerm)
	}
	if _, err := os.Stat(consts.OpenSolarIssuerDir); os.IsNotExist(err) {
		os.MkdirAll(consts.OpenSolarIssuerDir, os.ModePerm)
	}
}

// don't lock since boltdb can only process one operation at a time. As the application
// grows bigger, this would be a major reason to search for a new db system
func OpenDB() (*bolt.DB, error) {
	// we need to check and create this directory if it doesn't exist
	db, err := bolt.Open(consts.DbDir+"/yol.db", 0600, nil) // store this in its ownd database
	if err != nil {
		log.Println("Couldn't open database, exiting!")
		return db, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(ProjectsBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(InvestorBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(RecipientBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(ContractorBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(UserBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(BondBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(CoopBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(InspectorBucket) // the projects bucket contains all our projects
		if err != nil {
			return err
		}
		return nil
	})
	return db, err
}

// DeleteProject deltes a given value corresponding to the ky from the database
// DeleteProject should be used only in cases where something is wrong from our side
// while creating an project. For other cases, we should set Live to False and edit
// the project
// DeleteKeyFromBucket deletes a given key from the bucket bucketName but doesn
// not shift indices of elements succeeding the deleted element's index
func DeleteKeyFromBucket(key int, bucketName []byte) error {
	// deleting project might be dangerous since that would mess with the other
	// functions, have it in here for now, don't do too much with it / fiox retrieve all
	// to handle this case
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		b.Delete(utils.ItoB(key))
		return nil
	})
}
