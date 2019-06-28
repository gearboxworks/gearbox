package vmbox

import (
	"encoding/json"
	"gearbox/eventbroker/eblog"
	"github.com/gearboxworks/go-status/only"
	"io/ioutil"
	"os"
	"path/filepath"
)


func (me *Vm) ReadConfig() error {

	var err error
	var data []byte

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		file := me.Entry.VmDir.AddFileToPath("%s.json", me.Entry.Name)
		err = file.FileExists()
		if err != nil {
			break
		}

		data, err = ioutil.ReadFile(file.String())
		if err != nil {
			break
		}

		err = json.Unmarshal(data, &me.Entry)
		if err != nil {
			break
		}

		err = me.VerifyConfig()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "VM config loaded OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) WriteConfig() error {

	var err error
	var data []byte
	var perm os.FileMode

	perm = 0664

	for range only.Once {

		data, err = json.Marshal(&me.Entry)
		if err != nil {
			break
		}

		_, err = me.Entry.VmDir.CreateIfNotExists()
		if err != nil {
			break
		}

		file := me.Entry.VmDir.AddFileToPath("%s.json", me.Entry.Name)

		tempfile, err := ioutil.TempFile(me.Entry.VmDir.String(), filepath.Base(file.String()))
		if err != nil {
			break
		}
		name := tempfile.Name()
		defer os.Remove(name)

		if err = tempfile.Close(); err != nil {
			break
		}

		if err = ioutil.WriteFile(name, data, perm); err != nil {
			break
		}

		if err = os.Chmod(name, perm); err != nil {
			break
		}

		if err = os.Rename(name, file.String()); err != nil {
			break
		}

		eblog.Debug(me.EntityId, "VM config loaded OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) VerifyConfig() error {

	var err error

	for range only.Once {

		err = me.Entry.EnsureNotNil()
		if err != nil {
			break
		}

		if me.Entry.Name == "" {
			err = me.EntityName.ProduceError("VM doesn't have a name")
			break
		}

		if me.Entry.Version == "" {
			me.Entry.Version = "latest"
		}

		if me.Entry.Console.Host == "" {
			me.Entry.Console.Host = "localhost"
		}

		if me.Entry.Console.Port == "" {
			me.Entry.Console.Port = "2023"
		}

		if me.Entry.Console.ReadWait == 0 {
			me.Entry.Console.ReadWait = DefaultConsoleReadWait
		}

		if me.Entry.Console.OkString == "" {
			me.Entry.Console.OkString = DefaultConsoleOkString
		}

		if me.Entry.Console.WaitDelay == 0 {
			me.Entry.Console.WaitDelay = DefaultConsoleWaitDelay
		}

		if me.Entry.Ssh.Host == "" {
			me.Entry.Ssh.Host = "localhost"
		}

		if me.Entry.Ssh.Port == "" {
			me.Entry.Ssh.Port = "2222"
		}

		if me.Entry.IconFile == nil {
			err = me.osPaths.UserConfigDir.AddFileToPath(IconLogoPng).FileExists()
			if err != nil {
				err = nil
				// Not really an error.
			} else {
				me.Entry.IconFile = me.osPaths.UserConfigDir.AddFileToPath(IconLogoPng)
			}
		}

		if me.Entry.VmDir == nil {
			me.Entry.VmDir = me.osPaths.UserConfigDir.AddToPath("vm")
		}
		_, err = me.Entry.VmDir.CreateIfNotExists()
		if err != nil {
			break
		}

		me.Entry.retryMax = DefaultRetries
		me.Entry.retryDelay = DefaultVmWaitTime


		eblog.Debug(me.EntityId, "VM config is OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) ConfigExists() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		file := me.Entry.VmDir.AddFileToPath("%s.json", me.Entry.Name)
		err = file.FileExists()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "VM config exists")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

