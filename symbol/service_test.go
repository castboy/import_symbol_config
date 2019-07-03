/*******************************************************************************
 * // Copyright AnchyTec Corp. All Rights Reserved.
 * // SPDX-License-Identifier: Apache-2.0
 * // Author: shaozhiming
 ******************************************************************************/

package symbol

import (
	"testing"
)

func TestGetSymbolService(t *testing.T) {
	ss := GetSymbolService()
	if ss == nil {
		t.Error("GetPriceService Faild!")
	}

	ss.Start()
}
