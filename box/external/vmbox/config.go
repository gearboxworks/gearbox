package vmbox

import (
	"encoding/json"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
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

		file := osdirs.AddFilef(me.Entry.VmDir, JsonFilePattern, me.Entry.Name)
		err = osdirs.CheckFileExists(file)
		if err != nil {
			break
		}

		data, err = ioutil.ReadFile(file)
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

		_, err = osdirs.CreateIfNotExists(me.Entry.VmDir)
		if err != nil {
			break
		}

		file := osdirs.AddFilef(me.Entry.VmDir, JsonFilePattern, me.Entry.Name)

		tempfile, err := ioutil.TempFile(me.Entry.VmDir, filepath.Base(file))
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

		if err = os.Rename(name, file); err != nil {
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
			err = msgs.MakeError(me.EntityName, "VM doesn't have a name")
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
			me.Entry.Ssh.Host = DefaultSshHost
		}

		if me.Entry.Ssh.Port == "" {
			me.Entry.Ssh.Port = DefaultSshPort
		}

		if me.Entry.IconFile == "" {
			fp := me.osPaths.AddFileToUserConfigDir(IconLogoPng)
			if osdirs.CheckFileExists(fp) == nil {
				me.Entry.IconFile = fp
			}
		}

		if me.Entry.VmDir == "" {
			me.Entry.VmDir = me.osPaths.AppendToUserConfigDir("vm")
		}
		_, err = osdirs.CreateIfNotExists(me.Entry.VmDir)
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

		file := osdirs.AddFilef(me.Entry.VmDir, JsonFilePattern, me.Entry.Name)
		err = osdirs.CheckFileExists(file)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "VM config exists")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

func (me *Vm) DestroyConfig() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		file := osdirs.AddFilef(me.Entry.VmDir, JsonFilePattern, me.Entry.Name)
		err = osdirs.FileDelete(file)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "VM config removed")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}
