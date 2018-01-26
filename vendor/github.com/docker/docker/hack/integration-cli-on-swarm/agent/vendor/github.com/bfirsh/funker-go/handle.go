package funker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
)

// Handle a Funker function.
func Handle(handler interface***REMOVED******REMOVED***) error ***REMOVED***
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()
	if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 1 || handlerType.NumOut() != 1 ***REMOVED***
		return fmt.Errorf("Handler must be a function with a single parameter and single return value.")
	***REMOVED***
	argsValue := reflect.New(handlerType.In(0))

	listener, err := net.Listen("tcp", ":9999")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	conn, err := listener.Accept()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// We close listener, because we only allow single request.
	// Note that TCP "backlog" cannot be used for that purpose.
	// http://www.perlmonks.org/?node_id=940662
	if err = listener.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***
	argsJSON, err := ioutil.ReadAll(conn)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = json.Unmarshal(argsJSON, argsValue.Interface())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	ret := handlerValue.Call([]reflect.Value***REMOVED***argsValue.Elem()***REMOVED***)[0].Interface()
	retJSON, err := json.Marshal(ret)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if _, err = conn.Write(retJSON); err != nil ***REMOVED***
		return err
	***REMOVED***

	return conn.Close()
***REMOVED***
