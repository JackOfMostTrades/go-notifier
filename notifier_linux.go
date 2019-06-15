// Copyright 2016 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package notifier

/*
// Generate flags with "pkg-config --cflags libnotify" and "pkg-config --libs libnotify"
#cgo CFLAGS: -pthread -I/usr/include/gdk-pixbuf-2.0 -I/usr/include/libpng16 -I/usr/include/glib-2.0 -I/usr/lib/x86_64-linux-gnu/glib-2.0/include
#cgo LDFLAGS: -lnotify -lgdk_pixbuf-2.0 -lgio-2.0 -lgobject-2.0 -lglib-2.0
#include <libnotify/notify.h>
*/
import "C"

import (
	"errors"
	"os"
	"unsafe"
)

type linuxNotifier struct{}

// NewNotifier constructs notifier for Linux
func NewNotifier() (Notifier, error) {
	if C.notify_is_initted() == 0 {
		appName := C.CString(os.Args[0])
		defer C.free(unsafe.Pointer(appName))
		res := C.notify_init(appName)
		if res == 0 {
			return nil, errors.New("Unable to run notify_init.")
		}
	}
	return &linuxNotifier{}, nil
}

// DeliverNotification sends a notification
func (n linuxNotifier) DeliverNotification(notification Notification) error {
	title := C.CString(notification.Title)
	defer C.free(unsafe.Pointer(title))
	body := C.CString(notification.Message)
	defer C.free(unsafe.Pointer(body))
	imagePath := C.CString(notification.ImagePath)
	defer C.free(unsafe.Pointer(imagePath))

	note := C.notify_notification_new(title, body, imagePath)
	defer C.g_object_unref(C.gpointer(note))

	var gerr *C.GError
	res := C.notify_notification_show(note, &gerr)
	if res == 0 {
		msg := C.GoString(gerr.message)
		C.g_error_free(gerr)
		return errors.New(msg)
	}

	return nil
}
