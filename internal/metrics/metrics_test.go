/*
 * Copyright (C) 2023 akquinet GmbH
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metrics

import (
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTotalPostNotifyRequests(t *testing.T) {

	// Use testutil to get the current value of the metric with the specified label
	val := testutil.CollectAndCount(TotalPostNotifyRequests, "matrix_go_push_processed_post_notify_requests_total")
	// Assert that the metric has the expected value for the specified label
	assert.Equal(t, 1, int(val))
}
func TestTotalFirebaseRequestsMetric(t *testing.T) {
	// Use testutil to get the current value of the metric with the specified label
	val := testutil.CollectAndCount(TotalFirebaseRequestsMetric, "matrix_go_push_processed_firebase_requests_total")
	// Assert that the metric has the expected value for the specified label
	assert.Equal(t, 1, int(val))
}
func TestTotalFirbaseSentRequests(t *testing.T) {
	// Use testutil to get the current value of the metric with the specified label
	val := testutil.CollectAndCount(TotalFirbaseSentRequests, "matrix_go_push_processed_firebase_sent_requests_total")
	// Assert that the metric has the expected value for the specified label
	assert.Equal(t, 1, int(val))
}
func TestTotalFirebaseRejectedRequestsMetric(t *testing.T) {
	// Use testutil to get the current value of the metric with the specified label
	val := testutil.CollectAndCount(TotalFirebaseRejectedRequestsMetric, "matrix_go_push_processed_firebase_rejected_requests_total")
	// Assert that the metric has the expected value for the specified label
	assert.Equal(t, 1, int(val))
}
func TestTotalWrongDeviceIdsProcessedMetric(t *testing.T) {
	// Use testutil to get the current value of the metric with the specified label
	val := testutil.CollectAndCount(TotalWrongDeviceIdsProcessedMetric, "matrix_go_push_wrong_device_ids_sent_total")
	// Assert that the metric has the expected value for the specified label
	assert.Equal(t, 1, int(val))
}
func TestTotalMatrixNotFoundRequestsMetric(t *testing.T) {
	// Use testutil to get the current value of the metric with the specified label
	val := testutil.CollectAndCount(TotalMatrixNotFoundRequestsMetric, "matrix_go_push_matrix_not_found_total")
	// Assert that the metric has the expected value for the specified label
	assert.Equal(t, 1, int(val))
}
func TestTotalMatrixNotAllowedRequestsMetric(t *testing.T) {
	// Use testutil to get the current value of the metric with the specified label
	val := testutil.CollectAndCount(TotalMatrixNotAllowedRequestsMetric, "matrix_go_push_matrix_not_allowed_total")
	// Assert that the metric has the expected value for the specified label
	assert.Equal(t, 1, int(val))
}
