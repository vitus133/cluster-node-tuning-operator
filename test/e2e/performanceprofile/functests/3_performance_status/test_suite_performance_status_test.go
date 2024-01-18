//go:build !unittests
// +build !unittests

package __performance_status_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-logr/stdr"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/reporters"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/errors"

	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	testutils "github.com/openshift/cluster-node-tuning-operator/test/e2e/performanceprofile/functests/utils"
	testclient "github.com/openshift/cluster-node-tuning-operator/test/e2e/performanceprofile/functests/utils/client"
	testlog "github.com/openshift/cluster-node-tuning-operator/test/e2e/performanceprofile/functests/utils/log"
	"github.com/openshift/cluster-node-tuning-operator/test/e2e/performanceprofile/functests/utils/namespaces"

	qe_reporters "kubevirt.io/qe-tools/pkg/ginkgo-reporters"
)

var _ = BeforeSuite(func() {
	// create test namespace
	err := testclient.Client.Create(context.TODO(), namespaces.TestingNamespace)
	if errors.IsAlreadyExists(err) {
		testlog.Warning("test namespace already exists, that is unexpected")
		return
	}
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := testclient.Client.Delete(context.TODO(), namespaces.TestingNamespace)
	Expect(err).ToNot(HaveOccurred())
	err = namespaces.WaitForDeletion(testutils.NamespaceTesting, 5*time.Minute)
})

func TestPerformanceUpdate(t *testing.T) {
	ctrllog.SetLogger(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)))

	RegisterFailHandler(Fail)

	RunSpecs(t, "Performance Addon Operator configuration")
}

var _ = ReportAfterSuite("e2e serial suite", func(r Report) {
	if qe_reporters.Polarion.Run {
		reporters.ReportViaDeprecatedReporter(&qe_reporters.Polarion, r)
	}
})
