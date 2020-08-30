/*


Copyright (C) 2020  The XRootD Collaboration

This library is free software; you can redistribute it and/or
modify it under the terms of the GNU Lesser General Public
License as published by the Free Software Foundation; either
version 2.1 of the License, or (at your option) any later version.

This library is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public
License along with this library; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301
USA
*/

package watch

import "testing"

func TestMultipleWatchers(t *testing.T) {
	var watchers = Watchers([]Watcher{})
	if len(watchers.ToSlice()) != 0 {
		t.Errorf("watchers must be empty: %v", watchers)
	}
	watchers = watchers.AddWatchers(newTestWatcher(0))
	if len(watchers.ToSlice()) != 1 {
		t.Errorf("watchers must have 1 watcher: %v", watchers)
	}
}
