/*
Copyright © 2021 John Hooks john@hooks.technology

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

//import (
//	"fmt"
//	"os"
//	"time"
//
//	"github.com/fauna/faunadb-go/v3/faunadb"
//)
//
//type sizeError struct {
//	Details string
//}
//
//type Fauna struct {
//	conn       *faunadb.FaunaClient
//	collection string
//}
//
//func NewFaunaConnection(token, coll string) *Fauna {
//	c := faunadb.NewFaunaClient(token)
//
//	return &Fauna{
//		conn:       c,
//		collection: coll,
//	}
//}
//
//// Close exists to satisfy the client interface
//func (f Fauna) Close() {}
//
//// fiveDays is just a small function to get the timestamp
//// five days in the future.
//func fiveDays() time.Time {
//	now := time.Now()
//	fiveDays := 120 * time.Hour
//	return now.Add(fiveDays)
//}
//
//// BuildRecord does the initial record creation by generating a password
//// and encrypting the text and then storing that in the record.
//func (f *Fauna) Write(s Secret) error {
//
//	return f.write(s)
//}
//
//// AddRecord adds a document to the specified collection using
//// the data in the FaunaRecord receiver.
//func (f *Fauna) write(s Secret) error {
//	var ref faunadb.RefV
//
//	res, err := f.conn.Query(
//		faunadb.Create(
//			faunadb.Collection(f.collection),
//			faunadb.Obj{
//				"data": faunadb.Obj{
//					"id":    s.ID,
//					"text":  s.Text,
//					"views": s.Views,
//					"ttl":   fiveDays(),
//				}}))
//	if err != nil {
//		return fmt.Errorf("write: %w", err)
//	}
//
//	if err := res.At(faunadb.ObjKey("ref")).Get(&ref); err != nil {
//		return fmt.Errorf("write: %w", err)
//	}
//
//	return nil
//
//}
//
//// Read is just a wrapper to be able to check the number of views
//// left on the record. If the value is 0 it calls Delete.
//func (f *Fauna) Read(id string) (Secret, error) {
//	record := Secret{
//		ID: id,
//	}
//
//	if err := f.read(&record); err != nil {
//		return Secret{}, err
//	}
//
//	return record, nil
//
//}
//
//// readRecord retrieves a document from a collection in FaunaDB. It assigns the data
//// to the Record.
//func (f *Fauna) read(s *Secret) error {
//	resp, err := f.conn.Query(
//		faunadb.Get(faunadb.Ref(faunadb.Collection(f.collection), s.ID)))
//	if err != nil {
//		// check if error is a 404 or not
//		_, ok := err.(faunadb.NotFound)
//		if ok {
//			return NewSecretError(404, "message not found")
//		}
//		if !ok {
//			return fmt.Errorf("error getting record: %w", err)
//		}
//	}
//
//	if err := resp.At(faunadb.ObjKey("data")).Get(s); err != nil {
//		return fmt.Errorf("read: %w", err)
//	}
//
//	return nil
//}
//
//// Delete deletes a document from a collection in FaunaDB based on the ID
//// from the Record.
//func (f *Fauna) Delete(id string) error {
//	_, err := f.conn.Query(
//		faunadb.Delete(faunadb.Ref(faunadb.Collection(os.Getenv("FAUNA_COLLECTION")), id)))
//	if err != nil {
//		return fmt.Errorf("Delete: %w", err)
//	}
//
//	return nil
//}
