package maniladriver

import (
	"context"

	"github.com/go-logr/logr"
	maniladriverv1alpha1 "github.com/openshift/csi-driver-manila-operator/pkg/apis/maniladriver/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	labelsManilaControllerPlugin = map[string]string{
		"app":       "openstack-manila-csi",
		"component": "controllerplugin",
	}
)

func (r *ReconcileManilaDriver) handleManilaControllerPluginRBAC(instance *maniladriverv1alpha1.ManilaDriver, reqLogger logr.Logger) error {
	reqLogger.Info("Reconciling Manila Controller Plugin RBAC resources")

	// Manila Controller Plugin Service Account
	err := r.handleManilaControllerPluginServiceAccount(instance, reqLogger)
	if err != nil {
		return err
	}

	// Manila Controller Plugin Cluster Role
	err = r.handleManilaControllerPluginClusterRole(instance, reqLogger)
	if err != nil {
		return err
	}

	// Manila Controller Plugin Cluster Role Binding
	err = r.handleManilaControllerPluginClusterRoleBinding(instance, reqLogger)
	if err != nil {
		return err
	}

	// Manila Controller Plugin Role
	err = r.handleManilaControllerPluginRole(instance, reqLogger)
	if err != nil {
		return err
	}

	// Manila Controller Plugin Role Binding
	err = r.handleManilaControllerPluginRoleBinding(instance, reqLogger)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileManilaDriver) handleManilaControllerPluginServiceAccount(instance *maniladriverv1alpha1.ManilaDriver, reqLogger logr.Logger) error {
	reqLogger.Info("Reconciling Manila Controller Plugin Service Account")

	// Define a new ServiceAccount object
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openstack-manila-csi-controllerplugin",
			Namespace: "openshift-manila-csi-driver",
			Labels:    labelsManilaControllerPlugin,
		},
	}

	if err := annotator.SetLastAppliedAnnotation(sa); err != nil {
		return err
	}

	// Check if this ServiceAccount already exists
	found := &corev1.ServiceAccount{}
	err := r.apiReader.Get(context.TODO(), types.NamespacedName{Name: sa.Name, Namespace: sa.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ServiceAccount", "ServiceAccount.Namespace", sa.Namespace, "ServiceAccount.Name", sa.Name)
		err = r.client.Create(context.TODO(), sa)
		if err != nil {
			return err
		}

		// ServiceAccount created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	// Check if we need to update the object
	equal, err := compareLastAppliedAnnotations(found, sa)
	if err != nil {
		return err
	}

	if !equal {
		reqLogger.Info("Updating ServiceAccount with new changes", "ServiceAccount.Namespace", found.Namespace, "ServiceAccount.Name", found.Name)
		err = r.client.Update(context.TODO(), sa)
		if err != nil {
			return err
		}
	} else {
		// ServiceAccount already exists - don't requeue
		reqLogger.Info("Skip reconcile: ServiceAccount already exists", "ServiceAccount.Namespace", found.Namespace, "ServiceAccount.Name", found.Name)
	}

	return nil
}

func (r *ReconcileManilaDriver) handleManilaControllerPluginClusterRole(instance *maniladriverv1alpha1.ManilaDriver, reqLogger logr.Logger) error {
	reqLogger.Info("Reconciling Manila Controller Plugin Cluster Role")

	// Define a new ClusterRole object
	cr := generateManilaControllerPluginClusterRole()

	if err := annotator.SetLastAppliedAnnotation(cr); err != nil {
		return err
	}

	// Check if this ClusterRole already exists
	found := &rbacv1.ClusterRole{}
	err := r.apiReader.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: ""}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ClusterRole", "ClusterRole.Name", cr.Name)
		err = r.client.Create(context.TODO(), cr)
		if err != nil {
			return err
		}

		// ClusterRole created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	// Check if we need to update the object
	equal, err := compareLastAppliedAnnotations(found, cr)
	if err != nil {
		return err
	}

	if !equal {
		reqLogger.Info("Updating ClusterRole with new changes", "ClusterRole.Name", found.Name)
		err = r.client.Update(context.TODO(), cr)
		if err != nil {
			return err
		}
	} else {
		// ClusterRole already exists - don't requeue
		reqLogger.Info("Skip reconcile: ClusterRole already exists", "ClusterRole.Name", found.Name)
	}

	return nil
}

func (r *ReconcileManilaDriver) handleManilaControllerPluginClusterRoleBinding(instance *maniladriverv1alpha1.ManilaDriver, reqLogger logr.Logger) error {
	reqLogger.Info("Reconciling Manila Controller Plugin Cluster Role Binding")

	// Define a new ClusterRoleBinding object
	crb := generateManilaControllerPluginClusterRoleBinding()

	if err := annotator.SetLastAppliedAnnotation(crb); err != nil {
		return err
	}

	// Check if this ClusterRoleBinding already exists
	found := &rbacv1.ClusterRoleBinding{}
	err := r.apiReader.Get(context.TODO(), types.NamespacedName{Name: crb.Name, Namespace: ""}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ClusterRoleBinding", "ClusterRoleBinding.Name", crb.Name)
		err = r.client.Create(context.TODO(), crb)
		if err != nil {
			return err
		}

		// ClusterRoleBinding created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	// Check if we need to update the object
	equal, err := compareLastAppliedAnnotations(found, crb)
	if err != nil {
		return err
	}

	if !equal {
		reqLogger.Info("Updating ClusterRoleBinding with new changes", "ClusterRoleBinding.Name", found.Name)
		err = r.client.Update(context.TODO(), crb)
		if err != nil {
			return err
		}
	} else {
		// ClusterRoleBinding already exists - don't requeue
		reqLogger.Info("Skip reconcile: ClusterRoleBinding already exists", "ClusterRoleBinding.Name", found.Name)
	}

	return nil
}

func (r *ReconcileManilaDriver) handleManilaControllerPluginRole(instance *maniladriverv1alpha1.ManilaDriver, reqLogger logr.Logger) error {
	reqLogger.Info("Reconciling Manila Controller Plugin Role")

	// Define a new Role object
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openstack-manila-csi-controllerplugin",
			Namespace: "openshift-manila-csi-driver",
			Labels:    labelsManilaControllerPlugin,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
				Verbs:     []string{"get", "watch", "list", "delete", "update", "create"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get", "list", "watch", "create", "delete"},
			},
		},
	}

	if err := annotator.SetLastAppliedAnnotation(role); err != nil {
		return err
	}

	// Check if this Role already exists
	found := &rbacv1.Role{}
	err := r.apiReader.Get(context.TODO(), types.NamespacedName{Name: role.Name, Namespace: role.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Role", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
		err = r.client.Create(context.TODO(), role)
		if err != nil {
			return err
		}

		// Role created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	// Check if we need to update the object
	equal, err := compareLastAppliedAnnotations(found, role)
	if err != nil {
		return err
	}

	if !equal {
		reqLogger.Info("Updating Role with new changes", "Role.Namespace", found.Namespace, "Role.Name", found.Name)
		err = r.client.Update(context.TODO(), role)
		if err != nil {
			return err
		}
	} else {
		// Role already exists - don't requeue
		reqLogger.Info("Skip reconcile: Role already exists", "Role.Namespace", found.Namespace, "Role.Name", found.Name)
	}

	return nil
}

func (r *ReconcileManilaDriver) handleManilaControllerPluginRoleBinding(instance *maniladriverv1alpha1.ManilaDriver, reqLogger logr.Logger) error {
	reqLogger.Info("Reconciling Manila Controller Plugin Role Binding")

	// Define a new RoleBinding object
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openstack-manila-csi-controllerplugin",
			Namespace: "openshift-manila-csi-driver",
			Labels:    labelsManilaControllerPlugin,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "openstack-manila-csi-controllerplugin",
				Namespace: "openshift-manila-csi-driver",
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     "openstack-manila-csi-controllerplugin",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	if err := annotator.SetLastAppliedAnnotation(rb); err != nil {
		return err
	}

	// Check if this RoleBinding already exists
	found := &rbacv1.RoleBinding{}
	err := r.apiReader.Get(context.TODO(), types.NamespacedName{Name: rb.Name, Namespace: rb.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new RoleBinding", "RoleBinding.Namespace", rb.Namespace, "RoleBinding.Name", rb.Name)
		err = r.client.Create(context.TODO(), rb)
		if err != nil {
			return err
		}

		// RoleBinding created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	// Check if we need to update the object
	equal, err := compareLastAppliedAnnotations(found, rb)
	if err != nil {
		return err
	}

	if !equal {
		reqLogger.Info("Updating RoleBinding with new changes", "RoleBinding.Namespace", found.Namespace, "RoleBinding.Name", found.Name)
		err = r.client.Update(context.TODO(), rb)
		if err != nil {
			return err
		}
	} else {
		// RoleBinding already exists - don't requeue
		reqLogger.Info("Skip reconcile: RoleBinding already exists", "RoleBinding.Namespace", found.Namespace, "RoleBinding.Name", found.Name)
	}

	return nil
}

func (r *ReconcileManilaDriver) deleteManilaControllerPluginClusterRole(reqLogger logr.Logger) error {
	cr := generateManilaControllerPluginClusterRole()

	err := r.client.Delete(context.TODO(), cr)
	if err != nil {
		return err
	}

	reqLogger.Info("Cluster Role was deleted succesfully", "ClusterRole.Name", cr.Name)

	return nil
}

func (r *ReconcileManilaDriver) deleteManilaControllerPluginClusterRoleBinding(reqLogger logr.Logger) error {
	crb := generateManilaControllerPluginClusterRoleBinding()

	err := r.client.Delete(context.TODO(), crb)
	if err != nil {
		return err
	}

	reqLogger.Info("Cluster Role Binding was deleted succesfully", "ClusterRoleBinding.Name", crb.Name)

	return nil
}

func generateManilaControllerPluginClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "openstack-manila-csi-controllerplugin",
			Labels: labelsManilaControllerPlugin,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get", "list"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumes"},
				Verbs:     []string{"get", "list", "watch", "create", "delete"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumeclaims"},
				Verbs:     []string{"get", "list", "watch", "update"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"list", "watch", "create", "update", "patch"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"csinodes"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"snapshot.storage.k8s.io"},
				Resources: []string{"volumesnapshotclasses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"snapshot.storage.k8s.io"},
				Resources: []string{"volumesnapshotcontents"},
				Verbs:     []string{"create", "get", "list", "watch", "update", "delete"},
			},
			{
				APIGroups: []string{"snapshot.storage.k8s.io"},
				Resources: []string{"volumesnapshotcontents/status"},
				Verbs:     []string{"update"},
			},
			{
				APIGroups: []string{"snapshot.storage.k8s.io"},
				Resources: []string{"volumesnapshots"},
				Verbs:     []string{"get", "list", "watch", "update"},
			},
			{
				APIGroups: []string{"snapshot.storage.k8s.io"},
				Resources: []string{"volumesnapshots/status"},
				Verbs:     []string{"update"},
			},
			{
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"create", "list", "watch", "delete", "get", "update"},
			},
		},
	}
}

func generateManilaControllerPluginClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "openstack-manila-csi-controllerplugin",
			Labels: labelsManilaControllerPlugin,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "openstack-manila-csi-controllerplugin",
				Namespace: "openshift-manila-csi-driver",
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "openstack-manila-csi-controllerplugin",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}
