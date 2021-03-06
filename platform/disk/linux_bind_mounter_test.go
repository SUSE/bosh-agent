package disk_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-agent/platform/disk"
	fakedisk "github.com/cloudfoundry/bosh-agent/platform/disk/fakes"
)

var _ = Describe("linuxBindMounter", func() {
	var (
		delegateErr     error
		delegateMounter *fakedisk.FakeMounter
		mounter         Mounter
	)

	BeforeEach(func() {
		delegateErr = errors.New("fake-err")
		delegateMounter = &fakedisk.FakeMounter{}
		mounter = NewLinuxBindMounter(delegateMounter)
	})

	Describe("MountFilesystem", func() {
		Context("when mounting regular directory", func() {
			It("delegates to mounter and adds 'bind' option to mount as a bind-mount", func() {
				delegateMounter.MountFilesystemErr = delegateErr

				err := mounter.MountFilesystem("fake-partition-path", "fake-mount-path", "awesomefs", "fake-opt1")

				// Outputs
				Expect(err).To(Equal(delegateErr))

				// Inputs
				Expect(delegateMounter.MountFilesystemPartitionPaths).To(Equal([]string{"fake-partition-path"}))
				Expect(delegateMounter.MountFilesystemMountPoints).To(Equal([]string{"fake-mount-path"}))
				Expect(delegateMounter.MountFilesystemFstypes).To(Equal([]string{"awesomefs"}))
				Expect(delegateMounter.MountFilesystemMountOptions).To(Equal([][]string{{"fake-opt1", "bind"}}))
			})
		})

		Context("when mounting tmpfs", func() {
			It("delegates to mounter and does not add 'bind' option to mount as a bind-mount", func() {
				delegateMounter.MountFilesystemErr = delegateErr

				err := mounter.MountFilesystem("somesrc", "fake-mount-path", "tmpfs", "fake-opt1")

				// Outputs
				Expect(err).To(Equal(delegateErr))

				// Inputs
				Expect(delegateMounter.MountFilesystemPartitionPaths).To(Equal([]string{"somesrc"}))
				Expect(delegateMounter.MountFilesystemMountPoints).To(Equal([]string{"fake-mount-path"}))
				Expect(delegateMounter.MountFilesystemFstypes).To(Equal([]string{"tmpfs"}))
				Expect(delegateMounter.MountFilesystemMountOptions).To(Equal([][]string{{"fake-opt1"}}))
			})
		})
	})

	Describe("Mount", func() {
		It("delegates to mounter with empty string filesystem type (so that it can be inferred)", func() {
			delegateMounter.MountFilesystemErr = delegateErr

			err := mounter.Mount("fake-partition-path", "fake-mount-path", "fake-opt1")

			// Outputs
			Expect(err).To(Equal(delegateErr))

			// Inputs
			Expect(delegateMounter.MountFilesystemPartitionPaths).To(Equal([]string{"fake-partition-path"}))
			Expect(delegateMounter.MountFilesystemMountPoints).To(Equal([]string{"fake-mount-path"}))
			Expect(delegateMounter.MountFilesystemFstypes).To(Equal([]string{""}))
			Expect(delegateMounter.MountFilesystemMountOptions).To(Equal([][]string{{"fake-opt1", "bind"}}))
		})
	})

	Describe("RemountAsReadonly", func() {
		It("does not delegate to mounter because remount with 'bind' does not work", func() {
			err := mounter.RemountAsReadonly("fake-path")
			Expect(err).To(BeNil())
			Expect(delegateMounter.RemountAsReadonlyCalled).To(BeFalse())
		})
	})

	Describe("Remount", func() {
		It("delegates to mounter and adds 'bind' option to mount as a bind-mount", func() {
			delegateMounter.RemountErr = delegateErr

			err := mounter.Remount("fake-from-path", "fake-to-path", "fake-opt1")

			// Outputs
			Expect(err).To(Equal(delegateErr))

			// Inputs
			Expect(delegateMounter.RemountFromMountPoint).To(Equal("fake-from-path"))
			Expect(delegateMounter.RemountToMountPoint).To(Equal("fake-to-path"))
			Expect(delegateMounter.RemountMountOptions).To(Equal([]string{"fake-opt1", "bind"}))
		})
	})

	Describe("SwapOn", func() {
		It("delegates to mounter", func() {
			delegateMounter.SwapOnErr = delegateErr

			err := mounter.SwapOn("fake-path")

			// Outputs
			Expect(err).To(Equal(delegateErr))

			// Inputs
			Expect(delegateMounter.SwapOnPartitionPaths).To(Equal([]string{"fake-path"}))
		})
	})

	Describe("Unmount", func() {
		It("delegates to mounter", func() {
			delegateMounter.UnmountErr = delegateErr
			delegateMounter.UnmountDidUnmount = true

			didUnmount, err := mounter.Unmount("fake-device-path")

			// Outputs
			Expect(didUnmount).To(BeTrue())
			Expect(err).To(Equal(delegateErr))

			// Inputs
			Expect(delegateMounter.UnmountPartitionPathOrMountPoint).To(Equal("fake-device-path"))
		})
	})

	Describe("IsMountPoint", func() {
		It("delegates to mounter", func() {
			delegateMounter.IsMountPointErr = delegateErr
			delegateMounter.IsMountPointResult = true
			delegateMounter.IsMountPointPartitionPath = "fake-partition-path"

			partitionPath, isMountPoint, err := mounter.IsMountPoint("fake-device-path")

			// Outputs
			Expect(partitionPath).To(Equal("fake-partition-path"))
			Expect(isMountPoint).To(BeTrue())
			Expect(err).To(Equal(delegateErr))

			// Inputs
			Expect(delegateMounter.IsMountPointPath).To(Equal("fake-device-path"))
		})
	})

	Describe("IsMounted", func() {
		It("delegates to mounter", func() {
			delegateMounter.IsMountedErr = delegateErr
			delegateMounter.IsMountedResult = true

			isMounted, err := mounter.IsMounted("fake-device-path")

			// Outputs
			Expect(isMounted).To(BeTrue())
			Expect(err).To(Equal(delegateErr))

			// Inputs
			Expect(delegateMounter.IsMountedArgsForCall(0)).To(Equal("fake-device-path"))
		})
	})
})
