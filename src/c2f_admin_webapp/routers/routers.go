package routers

import (
	//	"C2FAdmin/controllers/customer"
	"c2f_admin_webapp/controllers/changePassword"

	"c2f_admin_webapp/controllers/analytics"
	"c2f_admin_webapp/controllers/analytics/clusterwise"
	"c2f_admin_webapp/controllers/analytics/forecastdashboard"
	"c2f_admin_webapp/controllers/clusterdashbord"
	"c2f_admin_webapp/controllers/clusterexecutivedashbord"
	"c2f_admin_webapp/controllers/customernotification"
	"c2f_admin_webapp/controllers/dashboard"
	"c2f_admin_webapp/controllers/notification"

	//	"c2f_admin_webapp/controllers/hohelpdeskdashbord"

	"c2f_admin_webapp/controllers/sysusers/createsysUser"
	"c2f_admin_webapp/controllers/sysusers/searchsysUser"
	"c2f_admin_webapp/controllers/sysusers/updatesysUser"
	"c2f_admin_webapp/controllers/sysusers/viewsysUser"

	"c2f_admin_webapp/controllers/collectionagent/createcollectionAgent"
	"c2f_admin_webapp/controllers/collectionagent/searchcollectionAgent"
	"c2f_admin_webapp/controllers/collectionagent/updatecollectionAgent"
	"c2f_admin_webapp/controllers/collectionagent/viewcollectionAgent"

	"c2f_admin_webapp/controllers/collectpayment/collectPayment"
	"c2f_admin_webapp/controllers/collectpayment/updatecollectPayment"

	"c2f_admin_webapp/controllers/users/createUser"
	"c2f_admin_webapp/controllers/users/searchUser"
	"c2f_admin_webapp/controllers/users/updateUser"
	"c2f_admin_webapp/controllers/users/viewUser"

	"c2f_admin_webapp/controllers/supplyexecutive/createSupplyExecutive"
	"c2f_admin_webapp/controllers/supplyexecutive/searchSupplyExecutive"
	"c2f_admin_webapp/controllers/supplyexecutive/updateSupplyExecutive"
	"c2f_admin_webapp/controllers/supplyexecutive/viewSupplyExecutive"

	"c2f_admin_webapp/controllers/clustersupplyexecutive/createClusterSupplyexecutive"
	"c2f_admin_webapp/controllers/clustersupplyexecutive/searchClusterSupplyexecutive"
	"c2f_admin_webapp/controllers/clustersupplyexecutive/updateClusterSupplyexecutive"
	"c2f_admin_webapp/controllers/clustersupplyexecutive/viewClusterSupplyexecutive"

	"c2f_admin_webapp/controllers/supplyorder/ordertosupplyexecutive"

	"c2f_admin_webapp/controllers/products/createProduct"
	"c2f_admin_webapp/controllers/products/searchProduct"
	"c2f_admin_webapp/controllers/products/updateProduct"
	"c2f_admin_webapp/controllers/products/updateProductPrice"

	"c2f_admin_webapp/controllers/products/viewProduct"

	"c2f_admin_webapp/controllers/clusterproducts/createClusterProduct"
	"c2f_admin_webapp/controllers/clusterproducts/searchClusterProduct"
	"c2f_admin_webapp/controllers/clusterproducts/updateClusterCharges"
	"c2f_admin_webapp/controllers/clusterproducts/updateClusterProduct"
	"c2f_admin_webapp/controllers/clusterproducts/viewClusterProduct"

	"c2f_admin_webapp/controllers/clusterusers/createClusterUser"
	"c2f_admin_webapp/controllers/clusterusers/searchClusterUser"
	"c2f_admin_webapp/controllers/clusterusers/updateClusterUser"
	"c2f_admin_webapp/controllers/clusterusers/viewClusterUser"

	"c2f_admin_webapp/controllers/clusterclusterexecutive/createClusterClusterexecutive"
	"c2f_admin_webapp/controllers/clusterclusterexecutive/searchClusterClusterexecutive"
	"c2f_admin_webapp/controllers/clusterclusterexecutive/updateClusterClusterexecutive"
	"c2f_admin_webapp/controllers/clusterclusterexecutive/viewClusterClusterexecutive"

	"c2f_admin_webapp/controllers/clusterdelivery/createClusterDelivery"
	"c2f_admin_webapp/controllers/clusterdelivery/searchClusterDelivery"
	"c2f_admin_webapp/controllers/clusterdelivery/updateClusterDelivery"
	"c2f_admin_webapp/controllers/clusterdelivery/viewClusterDelivery"

	"c2f_admin_webapp/controllers/clustercollectionagent/createClusterCollectionagent"
	"c2f_admin_webapp/controllers/clustercollectionagent/searchClusterCollectionagent"
	"c2f_admin_webapp/controllers/clustercollectionagent/updateClusterCollectionagent"
	"c2f_admin_webapp/controllers/clustercollectionagent/viewClusterCollectionagent"

	"c2f_admin_webapp/controllers/reverseinventory/createReverseInventory"
	"c2f_admin_webapp/controllers/reverseinventory/searchReverseInventory"
	"c2f_admin_webapp/controllers/reverseinventory/updateReverseInventory"
	"c2f_admin_webapp/controllers/reverseinventory/viewReverseInventory"

	"c2f_admin_webapp/controllers/roles/createRole"
	"c2f_admin_webapp/controllers/roles/searchRole"
	"c2f_admin_webapp/controllers/roles/updateRole"
	"c2f_admin_webapp/controllers/roles/viewRole"

	"c2f_admin_webapp/controllers/supplyforecast/createSupplyForecast"
	"c2f_admin_webapp/controllers/supplyforecast/searchSupplyForecast"
	"c2f_admin_webapp/controllers/supplyforecast/selectSupplyForcast"
	"c2f_admin_webapp/controllers/supplyforecast/updateSupplyForecast"
	"c2f_admin_webapp/controllers/supplyforecast/viewSupplyForecast"

	"c2f_admin_webapp/controllers/supplyorder/createSupplyOrder"
	"c2f_admin_webapp/controllers/supplyorder/searchSupplyOrder"
	"c2f_admin_webapp/controllers/supplyorder/selectSupplyOrder"
	"c2f_admin_webapp/controllers/supplyorder/updateSupplyOrder"
	"c2f_admin_webapp/controllers/supplyorder/viewSupplyOrder"

	"c2f_admin_webapp/controllers/productproducercluster/createProductProducerCluster"
	"c2f_admin_webapp/controllers/productproducercluster/searchProductProducerCluster"
	"c2f_admin_webapp/controllers/productproducercluster/updateProductProducerCluster"
	"c2f_admin_webapp/controllers/productproducercluster/viewProductProducerCluster"

	"c2f_admin_webapp/controllers/deliveryusers/createDeliveryUser"
	"c2f_admin_webapp/controllers/deliveryusers/searchDeliveryUser"
	"c2f_admin_webapp/controllers/deliveryusers/updateDeliveryUser"
	"c2f_admin_webapp/controllers/deliveryusers/viewDeliveryUser"

	"c2f_admin_webapp/controllers/clusters/createCluster"
	"c2f_admin_webapp/controllers/clusters/searchCluster"
	"c2f_admin_webapp/controllers/clusters/updateCluster"
	"c2f_admin_webapp/controllers/clusters/updateClusterProductPrices"
	"c2f_admin_webapp/controllers/clusters/viewCluster"

	"c2f_admin_webapp/controllers/customeroder/createOrder"
	"c2f_admin_webapp/controllers/customeroder/placeOrder"
	"c2f_admin_webapp/controllers/customeroder/searchOrder"
	"c2f_admin_webapp/controllers/customeroder/updateOrder"

	"c2f_admin_webapp/controllers/deliveryorder/createDeliverOder"
	"c2f_admin_webapp/controllers/deliveryorder/selectDeliverOder"
	"c2f_admin_webapp/controllers/deliveryorder/updateDeliverOder"
	"c2f_admin_webapp/controllers/deliveryorder/viewDeliverOder"

	"c2f_admin_webapp/controllers/purchaseorder/createPurchaseOrder"
	"c2f_admin_webapp/controllers/purchaseorder/selectPurchaseOrder"
	"c2f_admin_webapp/controllers/purchaseorder/updatePurchaseOrder"
	"c2f_admin_webapp/controllers/purchaseorder/viewPurchaseOrder"

	"c2f_admin_webapp/controllers/supplyproductorder/createSupplyProductOrder"
	"c2f_admin_webapp/controllers/supplyproductorder/selectSupplyProductOrder"
	"c2f_admin_webapp/controllers/supplyproductorder/updateSupplyProductOrder"
	"c2f_admin_webapp/controllers/supplyproductorder/viewSupplyProductOrder"

	// "c2f_admin_webapp/controllers/loadmap/createLoadMap"
	"c2f_admin_webapp/controllers/loadmap/selectLoadMap"
	// "c2f_admin_webapp/controllers/loadmap/updateLoadMap"
	"c2f_admin_webapp/controllers/loadmap/viewLoadMap"

	"c2f_admin_webapp/controllers/ordercollectionagent/createOrderCollectionagent"
	"c2f_admin_webapp/controllers/ordercollectionagent/selectOrderCollectionagent"
	"c2f_admin_webapp/controllers/ordercollectionagent/updateOrderCollectionagent"
	"c2f_admin_webapp/controllers/ordercollectionagent/viewOrderCollectionagent"

	"c2f_admin_webapp/controllers/producersupplyorder/searchProducerSupplyOrder"
	"c2f_admin_webapp/controllers/producersupplyorder/updateProducerSupplyOrder"
	"c2f_admin_webapp/controllers/producersupplyorder/viewProducerSupplyOrder"

	"c2f_admin_webapp/controllers/collectsupplyorder/searchCollectSupplyOrder"
	"c2f_admin_webapp/controllers/collectsupplyorder/updateCollectSupplyOrder"
	"c2f_admin_webapp/controllers/collectsupplyorder/viewCollectSupplyOrder"

	"c2f_admin_webapp/controllers/completeuspplyorder/searchCompleteSupplyOrder"
	"c2f_admin_webapp/controllers/completeuspplyorder/updateCompleteSupplyOrder"
	"c2f_admin_webapp/controllers/completeuspplyorder/viewCompleteSupplyOrder"

	"c2f_admin_webapp/controllers/reports/adminAduitReport"
	"c2f_admin_webapp/controllers/reports/clusterproductpurchaseReport"
	"c2f_admin_webapp/controllers/reports/collectionAgentReport"
	"c2f_admin_webapp/controllers/reports/customerorderReport"
	"c2f_admin_webapp/controllers/reports/deliveryorderReport"
	"c2f_admin_webapp/controllers/reports/inventoryReport"
	"c2f_admin_webapp/controllers/reports/supplyforcastReport"
	"c2f_admin_webapp/controllers/reports/supplyorderAssignmentReport"
	"c2f_admin_webapp/controllers/reports/supplyorderReport"

	"c2f_admin_webapp/controllers/invoice/clusterproductpurchaseinvoice"
	"c2f_admin_webapp/controllers/invoice/collectpaymentinvoice"
	"c2f_admin_webapp/controllers/invoice/customerbillinvoice"
	"c2f_admin_webapp/controllers/invoice/customerorderinvoice"
	"c2f_admin_webapp/controllers/invoice/deliveryorderinvoice"
	"c2f_admin_webapp/controllers/invoice/deliveryorderinvoice1"
	"c2f_admin_webapp/controllers/invoice/supplyorderinvoice"

	"c2f_admin_webapp/controllers/pushnotification"

	"c2f_admin_webapp/controllers/customer/createCustomer"
	"c2f_admin_webapp/controllers/customer/searchCustomer"
	"c2f_admin_webapp/controllers/customer/updateCustomer"
	"c2f_admin_webapp/controllers/customer/viewCustomer"

	"c2f_admin_webapp/controllers/category/createCategory"
	"c2f_admin_webapp/controllers/category/searchCategory"
	"c2f_admin_webapp/controllers/category/updateCategory"

	"c2f_admin_webapp/controllers/category/viewCategory"

	"c2f_admin_webapp/controllers/ownerinfo/ownerInfo"

	"c2f_admin_webapp/controllers/subcategory/createSubcategory"
	"c2f_admin_webapp/controllers/subcategory/searchSubcategory"
	"c2f_admin_webapp/controllers/subcategory/updateSubcategory"
	"c2f_admin_webapp/controllers/subcategory/viewSubcategory"

	"c2f_admin_webapp/controllers/completedelivery/searchDelivery"
	"c2f_admin_webapp/controllers/completedelivery/updateDelivery"
	"c2f_admin_webapp/controllers/completedelivery/viewDelivery"

	"c2f_admin_webapp/controllers/reassignorder/searchReassignOrder"
	"c2f_admin_webapp/controllers/reassignorder/updateReassignOrder"
	"c2f_admin_webapp/controllers/reassignorder/viewReassignOrder"

	"c2f_admin_webapp/controllers/forgotpassword"

	"c2f_admin_webapp/controllers/bills/customerorderBill"

	"c2f_admin_webapp/controllers/analytics/forecast/customerForecastReport"
	"c2f_admin_webapp/controllers/analytics/forecast/producerForecastReport"
	"c2f_admin_webapp/controllers/customer/customerBulkUpload"
	"c2f_admin_webapp/controllers/forecast/createCustomerForecast"
	"c2f_admin_webapp/controllers/forecast/createProducerForecast"
	"c2f_admin_webapp/controllers/users/producerBulkUpload"

	"c2f_admin_webapp/controllers/login"
	"c2f_admin_webapp/controllers/logout"

	"c2f_admin_webapp/controllers/bills/clusterproductpurchaseBill"
	"c2f_admin_webapp/controllers/bills/supplyorderBill"
	"c2f_admin_webapp/controllers/invoice/pickupinvoice"

	"c2f_admin_webapp/controllers/channelpartnertype/createChannelPartnerType"
	"c2f_admin_webapp/controllers/channelpartnertype/searchChannelPartnerType"
	"c2f_admin_webapp/controllers/channelpartnertype/updateChannelPartnerType"
	"c2f_admin_webapp/controllers/channelpartnertype/viewChannelPartnerType"

	"c2f_admin_webapp/controllers/channelpartner/createChannelPartner"
	"c2f_admin_webapp/controllers/channelpartner/searchChannelPartner"
	"c2f_admin_webapp/controllers/channelpartner/updateChannelPartner"
	"c2f_admin_webapp/controllers/channelpartner/viewChannelPartner"

	"c2f_admin_webapp/controllers/channelpartnercluster/createChannelPartnerCluster"
	"c2f_admin_webapp/controllers/channelpartnercluster/searchChannelPartnerCluster"
	"c2f_admin_webapp/controllers/channelpartnercluster/updateChannelPartnerCluster"
	"c2f_admin_webapp/controllers/channelpartnercluster/viewChannelPartnerCluster"

	"c2f_admin_webapp/controllers/channelpartnerproduct/createChannelPartnerProduct"
	"c2f_admin_webapp/controllers/channelpartnerproduct/searchChannelPartnerProduct"
	"c2f_admin_webapp/controllers/channelpartnerproduct/updateChannelPartnerProduct"
	"c2f_admin_webapp/controllers/channelpartnerproduct/viewChannelPartnerProduct"

	"c2f_admin_webapp/controllers/reports/surplusTrackingReport"
	"c2f_admin_webapp/controllers/reports/surplusmanagementReport"

	"c2f_admin_webapp/controllers/payupage"
	"c2f_admin_webapp/controllers/payuresponse"

	"c2f_admin_webapp/controllers/error"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/session"

	"github.com/astaxie/beego"
)

func init() {
	_ = session.Init()
	_ = db.Init()
	beego.ErrorController(&error.Error{})
	beego.Router(beego.AppConfig.String("MAIN_PATH"), &login.Login{})
	beego.Router(beego.AppConfig.String("LOGIN_PATH"), &login.Login{})
	beego.Router(beego.AppConfig.String("LOGOUT_PATH"), &logout.Logout{})
	beego.Router(beego.AppConfig.String("DASHBOARD_PATH"), &dashboard.Dashboard{})
	beego.Router(beego.AppConfig.String("CHANGE_PASSWORD_PATH"), &changePassword.ChangePassword{})
	beego.Router(beego.AppConfig.String("ANALYTICS_PATH"), &analytics.Analytics{})
	beego.Router(beego.AppConfig.String("CLUSTER_WISE_DATA_PATH"), &clusterwise.ClusterwiseData{})
	beego.Router(beego.AppConfig.String("NOTIFICATION_PATH"), &notification.Notification{})
	beego.Router(beego.AppConfig.String("CUSTOMER_NOTIFICATION_PATH"), &customernotification.CustomerNotification{})

	beego.Router(beego.AppConfig.String("CREATE_SYS_USER_PATH"), &createsysUser.CreatesysUser{})
	beego.Router(beego.AppConfig.String("SEACRH_SYS_USER_PATH"), &searchsysUser.SearchsysUser{})
	beego.Router(beego.AppConfig.String("UPDATE_SYS_USER_PATH"), &updatesysUser.UpdatesysUser{})
	beego.Router(beego.AppConfig.String("VIEW_SYS_USER_PATH"), &viewsysUser.ViewsysUser{})

	beego.Router(beego.AppConfig.String("CREATE_COLLECTION_AGENT_PATH"), &createcollectionAgent.CreatecollectionAgent{})
	beego.Router(beego.AppConfig.String("SEACRH_COLLECTION_AGENT_PATH"), &searchcollectionAgent.SearchcollectionAgent{})
	beego.Router(beego.AppConfig.String("UPDATE_COLLECTION_AGENT_PATH"), &updatecollectionAgent.UpdatecollectionAgent{})
	beego.Router(beego.AppConfig.String("VIEW_COLLECTION_AGENT_PATH"), &viewcollectionAgent.ViewcollectionAgent{})

	beego.Router(beego.AppConfig.String("COLLECT_PAYMENT_PATH"), &collectPayment.CollectPayment{})
	beego.Router(beego.AppConfig.String("UPDATE_COLLECT_PAYMENT_PATH"), &updatecollectPayment.UpdatecollectPayment{})

	beego.Router(beego.AppConfig.String("CREATE_USER_PATH"), &createUser.CreateUser{})
	beego.Router(beego.AppConfig.String("SEACRH_USER_PATH"), &searchUser.SearchUser{})
	beego.Router(beego.AppConfig.String("UPDATE_USER_PATH"), &updateUser.UpdateUser{})
	beego.Router(beego.AppConfig.String("VIEW_USER_PATH"), &viewUser.ViewUser{})

	beego.Router(beego.AppConfig.String("CREATE_PRODUCT_PATH"), &createProduct.CreateProduct{})
	beego.Router(beego.AppConfig.String("SEACRH_PRODUCT_PATH"), &searchProduct.SearchProduct{})
	beego.Router(beego.AppConfig.String("UPDATE_PRODUCT_PATH"), &updateProduct.UpdateProduct{})
	beego.Router(beego.AppConfig.String("VIEW_PRODUCT_PATH"), &viewProduct.ViewProduct{})
	beego.Router(beego.AppConfig.String("UPDATE_PRODUCT_PRICE_PATH"), &updateProductPrice.UpdateProductPrice{})

	beego.Router(beego.AppConfig.String("CREATE_CLUSTER_PRODUCT_PATH"), &createClusterProduct.CreateClusterProduct{})
	beego.Router(beego.AppConfig.String("SEACRH_CLUSTER_PRODUCT_PATH"), &searchClusterProduct.SearchClusterProduct{})
	beego.Router(beego.AppConfig.String("UPDATE_CLUSTER_PRODUCT_PATH"), &updateClusterProduct.UpdateClusterProduct{})
	beego.Router(beego.AppConfig.String("VIEW_CLUSTER_PRODUCT_PATH"), &viewClusterProduct.ViewClusterProduct{})
	beego.Router(beego.AppConfig.String("UPDATE_CLUSTER_CHARGES_PATH"), &updateClusterCharges.UpdateClusterCharges{})

	beego.Router(beego.AppConfig.String("CREATE_CLUSTER_USER_PATH"), &createClusterUser.CreateClusterUser{})
	beego.Router(beego.AppConfig.String("SEACRH_CLUSTER_USER_PATH"), &searchClusterUser.SearchClusterUser{})
	beego.Router(beego.AppConfig.String("UPDATE_CLUSTER_USER_PATH"), &updateClusterUser.UpdateClusterUser{})
	beego.Router(beego.AppConfig.String("VIEW_CLUSTER_USER_PATH"), &viewClusterUser.ViewClusterUser{})

	beego.Router(beego.AppConfig.String("CREATE_CLUSTER_CLUSTEREXECUTIVE_PATH"), &createClusterClusterexecutive.CreateClusterClusterexecutive{})
	beego.Router(beego.AppConfig.String("SEACRH_CLUSTER_CLUSTEREXECUTIVE_PATH"), &searchClusterClusterexecutive.SearchClusterClusterexecutive{})
	beego.Router(beego.AppConfig.String("UPDATE_CLUSTER_CLUSTEREXECUTIVE_PATH"), &updateClusterClusterexecutive.UpdateClusterClusterexecutive{})
	beego.Router(beego.AppConfig.String("VIEW_CLUSTER_CLUSTEREXECUTIVE_PATH"), &viewClusterClusterexecutive.ViewClusterClusterexecutive{})

	beego.Router(beego.AppConfig.String("CREATE_CLUSTER_DELIVERY_PATH"), &createClusterDelivery.CreateClusterDelivery{})
	beego.Router(beego.AppConfig.String("SEACRH_CLUSTER_DELIVERY_PATH"), &searchClusterDelivery.SearchClusterDelivery{})
	beego.Router(beego.AppConfig.String("UPDATE_CLUSTER_DELIVERY_PATH"), &updateClusterDelivery.UpdateClusterDelivery{})
	beego.Router(beego.AppConfig.String("VIEW_CLUSTER_DELIVERY_PATH"), &viewClusterDelivery.ViewClusterDelivery{})

	beego.Router(beego.AppConfig.String("CREATE_CLUSTER_COLLECTION_AGENT_PATH"), &createClusterCollectionagent.CreateClusterCollectionagent{})
	beego.Router(beego.AppConfig.String("SEACRH_CLUSTER_COLLECTION_AGENT_PATH"), &searchClusterCollectionagent.SearchClusterCollectionagent{})
	beego.Router(beego.AppConfig.String("UPDATE_CLUSTER_COLLECTION_AGENT_PATH"), &updateClusterCollectionagent.UpdateClusterCollectionagent{})
	beego.Router(beego.AppConfig.String("VIEW_CLUSTER_COLLECTION_AGENT_PATH"), &viewClusterCollectionagent.ViewClusterCollectionagent{})

	beego.Router(beego.AppConfig.String("CREATE_REVERSE_INVENTORY_PATH"), &createReverseInventory.CreateReverseInventory{})
	beego.Router(beego.AppConfig.String("SEACRH_REVERSE_INVENTORY_PATH"), &searchReverseInventory.SearchReverseInventory{})
	beego.Router(beego.AppConfig.String("UPDATE_REVERSE_INVENTORY_PATH"), &updateReverseInventory.UpdateReverseInventory{})
	beego.Router(beego.AppConfig.String("VIEW_REVERSE_INVENTORY_PATH"), &viewReverseInventory.ViewReverseInventory{})

	beego.Router(beego.AppConfig.String("CREATE_ROLE_PATH"), &createRole.CreateRole{})
	beego.Router(beego.AppConfig.String("SEACRH_ROLE_PATH"), &searchRole.SearchRole{})
	beego.Router(beego.AppConfig.String("UPDATE_ROLE_PATH"), &updateRole.UpdateRole{})
	beego.Router(beego.AppConfig.String("VIEW_ROLE_PATH"), &viewRole.ViewRole{})

	beego.Router(beego.AppConfig.String("CREATE_PRODUCT_PRODUCER_CLUSTER_PATH"), &createProductProducerCluster.CreateProductProducerCluster{})
	beego.Router(beego.AppConfig.String("SEACRH_PRODUCT_PRODUCER_CLUSTER_PATH"), &searchProductProducerCluster.SearchProductProducerCluster{})
	beego.Router(beego.AppConfig.String("UPDATE_PRODUCT_PRODUCER_CLUSTER_PATH"), &updateProductProducerCluster.UpdateProductProducerCluster{})
	beego.Router(beego.AppConfig.String("VIEW_PRODUCT_PRODUCER_CLUSTER_PATH"), &viewProductProducerCluster.ViewProductProducerCluster{})

	beego.Router(beego.AppConfig.String("CREATE_SUPPLYFORCAST_PATH"), &createSupplyForecast.CreateSupplyForecast{})
	beego.Router(beego.AppConfig.String("SELECT_SUPPLYFORCAST_PATH"), &selectSupplyForcast.SelectSupplyForcast{})
	beego.Router(beego.AppConfig.String("SEACRH_SUPPLYFORCAST_PATH"), &searchSupplyForecast.SearchSupplyForecast{})
	beego.Router(beego.AppConfig.String("UPDATE_SUPPLYFORCAST_PATH"), &updateSupplyForecast.UpdateSupplyForecast{})
	beego.Router(beego.AppConfig.String("VIEW_SUPPLYFORCAST_PATH"), &viewSupplyForecast.ViewSupplyForecast{})

	beego.Router(beego.AppConfig.String("CREATE_DELIVERY_USER_PATH"), &createDeliveryUser.CreateDeliveryUser{})
	beego.Router(beego.AppConfig.String("SEACRH_DELIVERY_USER_PATH"), &searchDeliveryUser.SearchDeliveryUser{})
	beego.Router(beego.AppConfig.String("UPDATE_DELIVERY_USER_PATH"), &updateDeliveryUser.UpdateDeliveryUser{})
	beego.Router(beego.AppConfig.String("VIEW_DELIVERY_USER_PATH"), &viewDeliveryUser.ViewDeliveryUser{})

	beego.Router(beego.AppConfig.String("CREATE_CLUSTER_PATH"), &createCluster.CreateCluster{})
	beego.Router(beego.AppConfig.String("SEACRH_CLUSTER_PATH"), &searchCluster.SearchCluster{})
	beego.Router(beego.AppConfig.String("UPDATE_CLUSTER_PATH"), &updateCluster.UpdateCluster{})
	beego.Router(beego.AppConfig.String("VIEW_CLUSTER_PATH"), &viewCluster.ViewCluster{})
	beego.Router(beego.AppConfig.String("UPDATE_CLUSTER_PRODUCT_PRICE_PATH"), &updateClusterProductPrices.UpdateClusterProductPrice{})

	beego.Router(beego.AppConfig.String("CREATE_SUPPLYORDER_PATH"), &createSupplyOrder.CreateSupplyOrder{})
	beego.Router(beego.AppConfig.String("SELECT_SUPPLYORDER_PATH"), &selectSupplyOrder.SelectSupplyOrder{})
	beego.Router(beego.AppConfig.String("SEACRH_SUPPLYORDER_PATH"), &searchSupplyOrder.SearchSupplyOrder{})
	beego.Router(beego.AppConfig.String("UPDATE_SUPPLYORDER_PATH"), &updateSupplyOrder.UpdateSupplyOrder{})
	beego.Router(beego.AppConfig.String("VIEW_SUPPLYORDER_PATH"), &viewSupplyOrder.ViewSupplyOrder{})

	beego.Router(beego.AppConfig.String("SEACRH_PRODUCER_SUPPLYORDER_PATH"), &searchProducerSupplyOrder.SearchProducerSupplyOrder{})
	beego.Router(beego.AppConfig.String("UPDATE_PRODUCER_SUPPLYORDER_PATH"), &updateProducerSupplyOrder.UpdateProducerSupplyOrder{})
	beego.Router(beego.AppConfig.String("VIEW_PRODUCER_SUPPLYORDER_PATH"), &viewProducerSupplyOrder.ViewProducerSupplyOrder{})

	beego.Router(beego.AppConfig.String("SUPPLY_FORCAST_REPORT_PATH"), &supplyforcastReport.SupplyforcastReport{})
	beego.Router(beego.AppConfig.String("SUPPLY_ORDER_REPORT_PATH"), &supplyorderReport.SupplyorderReport{})
	beego.Router(beego.AppConfig.String("INVERNTORY_REPORT_PATH"), &inventoryReport.InventoryReport{})
	beego.Router(beego.AppConfig.String("CUSTOMER_ORDER_REPORT_PATH"), &customerorderReport.CustomerorderReport{})
	beego.Router(beego.AppConfig.String("DELIVERY_ORDER_REPORT_PATH"), &deliveryorderReport.DeliveryorderReport{})
	beego.Router(beego.AppConfig.String("ADUIT_ORDER_REPORT_PATH"), &adminAduitReport.AdminAduitReport{})
	beego.Router(beego.AppConfig.String("COLLECTION_AGENT_REPORT_PATH"), &collectionAgentReport.CollectionAgentReport{})
	beego.Router(beego.AppConfig.String("CLUSTER_PRODUCT_PURCHASE_REPORT_PATH"), &clusterproductpurchaseReport.ClusterproductpurchaseReport{})
	beego.Router(beego.AppConfig.String("SUPPLY_ORDER_ASSIGNMENT_REPORT_PATH"), &supplyorderAssignmentReport.SupplyorderAssignmentReport{})

	beego.Router(beego.AppConfig.String("DELIVERY_ORDER_INVOICE_PATH"), &deliveryorderinvoice.Deliveryorderinvoice{})
	beego.Router(beego.AppConfig.String("DELIVERY_ORDER_INVOICE_PATH1"), &deliveryorderinvoice1.Deliveryorderinvoice1{})

	beego.Router(beego.AppConfig.String("SEACRH_COLLECT_SUPPLYORDER_PATH"), &searchCollectSupplyOrder.SearchCollectSupplyOrder{})
	beego.Router(beego.AppConfig.String("UPDATE_COLLECT_SUPPLYORDER_PATH"), &updateCollectSupplyOrder.UpdateCollectSupplyOrder{})
	beego.Router(beego.AppConfig.String("VIEW_COLLECT_SUPPLYORDER_PATH"), &viewCollectSupplyOrder.ViewCollectSupplyOrder{})

	beego.Router(beego.AppConfig.String("SEACRH_COMPLETE_SUPPLYDEMAND_PATH"), &searchCompleteSupplyOrder.SearchCompleteSupplyOrder{})
	beego.Router(beego.AppConfig.String("UPDATE_COMPLETE_SUPPLYDEMAND_PATH"), &updateCompleteSupplyOrder.UpdateCompleteSupplyOrder{})
	beego.Router(beego.AppConfig.String("VIEW_COMPLETE_SUPPLYDEMAND_PATH"), &viewCompleteSupplyOrder.ViewCompleteSupplyOrder{})

	beego.Router(beego.AppConfig.String("CREATE_CUSTOMER_ORDER_PATH"), &createOrder.CreateOrder{})
	beego.Router(beego.AppConfig.String("SEACRH_CUSTOMER_ORDER_PATH"), &searchOrder.SearchOrder{})
	beego.Router(beego.AppConfig.String("UPDATE_CUSTOMER_ORDER_PATH"), &updateOrder.UpdateOrder{})
	beego.Router(beego.AppConfig.String("PLACE_CUSTOMER_ORDER_PATH"), &placeOrder.PlaceOrder{})
	beego.Router(beego.AppConfig.String("CUSTOMER_ADDITIONAL_ADDRESS_PATH"), &placeOrder.GetCustomerAdditionalAddress{})

	beego.Router(beego.AppConfig.String("CREATE_CUSTOMER_PATH"), &createCustomer.CreateCustomer{})
	beego.Router(beego.AppConfig.String("SEACRH_CUSTOMER_PATH"), &searchCustomer.SearchCustomer{})
	beego.Router(beego.AppConfig.String("UPDATE_CUSTOMER_PATH"), &updateCustomer.UpdateCustomer{})
	beego.Router(beego.AppConfig.String("VIEW_CUSTOMER_PATH"), &viewCustomer.ViewCustomer{})

	beego.Router(beego.AppConfig.String("CREATE_PURCHASE_ORDER_PATH"), &createPurchaseOrder.CreatePurchaseOrder{})
	beego.Router(beego.AppConfig.String("SELECT_PURCHASE_ORDER_PATH"), &selectPurchaseOrder.SelectPurchaseOrder{})
	beego.Router(beego.AppConfig.String("UPDATE_PURCHASE_ORDER_PATH"), &updatePurchaseOrder.UpdatePurchaseOrder{})
	beego.Router(beego.AppConfig.String("VIEW_PURCHASE_ORDER_PATH"), &viewPurchaseOrder.ViewPurchaseOrder{})

	beego.Router(beego.AppConfig.String("CREATE_SUPPLY_PRODUCT_ORDER_PATH"), &createSupplyProductOrder.CreateSupplyProductOrder{})
	beego.Router(beego.AppConfig.String("SELECT_SUPPLY_PRODUCT_ORDER_PATH"), &selectSupplyProductOrder.SelectSupplyProductOrder{})
	beego.Router(beego.AppConfig.String("UPDATE_SUPPLY_PRODUCT_ORDER_PATH"), &updateSupplyProductOrder.UpdateSupplyProductOrder{})
	beego.Router(beego.AppConfig.String("VIEW_SUPPLY_PRODUCT_ORDER_PATH"), &viewSupplyProductOrder.ViewSupplyProductOrder{})

	beego.Router(beego.AppConfig.String("CREATE_DELIVERYUSER_ODER_PATH"), &createDeliveryOder.CreateDeliveryOder{})
	beego.Router(beego.AppConfig.String("SELECT_DELIVERYUSER_ODER_PATH"), &selectDeliveryOder.SelectDeliveryOder{})
	beego.Router(beego.AppConfig.String("GET_DELIVERYDATE_ORDERS_PATH"), &selectDeliveryOder.GetDeliveryDateOrders{})
	beego.Router(beego.AppConfig.String("UPDATE_DELIVERYUSER_ODER_PATH"), &updateDeliveryOder.UpdateDeliveryOder{})
	beego.Router(beego.AppConfig.String("VIEW_DELIVERYUSER_ODER_PATH"), &viewDeliveryOder.ViewDeliveryOder{})

	// beego.Router(beego.AppConfig.String("CREATE_LOAD_MAP_PATH"), &createLoadMap.CreateLoadMap{})
	beego.Router(beego.AppConfig.String("SELECT_LOAD_MAP_PATH"), &selectLoadMap.SelectLoadMap{})
	// beego.Router(beego.AppConfig.String("UPDATE_LOAD_MAP_PATH"), &updateLoadMap.UpdateLoadMap{})
	beego.Router(beego.AppConfig.String("VIEW_LOAD_MAP_PATH"), &viewLoadMap.ViewLoadMap{})

	beego.Router(beego.AppConfig.String("CREATE_ORDER_COLLECTIONAGENT_PATH"), &createOrderCollectionagent.CreateOrderCollectionagent{})
	beego.Router(beego.AppConfig.String("SELECT_ORDER_COLLECTIONAGENT_PATH"), &selectOrderCollectionagent.SelectOrderCollectionagent{})
	beego.Router(beego.AppConfig.String("UPDATE_ORDER_COLLECTIONAGENT_PATH"), &updateOrderCollectionagent.UpdateOrderCollectionagent{})
	beego.Router(beego.AppConfig.String("VIEW_ORDER_COLLECTIONAGENT_PATH"), &viewOrderCollectionagent.ViewOrderCollectionagent{})

	beego.Router(beego.AppConfig.String("PRODUCER_SUPPLYORDER_BILL_PATH"), &supplyorderBill.SupplyorderBill{})
	beego.Router(beego.AppConfig.String("SUPPLY_ORDER_INVOICE_PATH"), &supplyorderinvoice.SupplyorderInvoice{})
	beego.Router(beego.AppConfig.String("CUSTOMER_ORDER_INVOICE_PATH"), &customerorderinvoice.CustomerorderInvoice{})
	beego.Router(beego.AppConfig.String("CLUSTER_PRODUCT_PURCHASE_INVOICE_PATH"), &clusterproductpurchaseinvoice.ClusterproductPurchaseInvoice{})
	beego.Router(beego.AppConfig.String("CUSTOMER_BILL_INVOICE_PATH"), &customerbillinvoice.CustomerBillInvoice{})
	beego.Router(beego.AppConfig.String("COLLECT_PAYMENT_INVOICE_PATH"), &collectpaymentinvoice.CollectPaymentInvoice{})

	beego.Router(beego.AppConfig.String("PUSH_NOTIFICATION"), &pushNotification.PushNotification{})

	beego.Router(beego.AppConfig.String("FORGOT_PASSWORD"), &forgotpassword.Forgotpassword{})
	beego.Router(beego.AppConfig.String("OWNER_INFO_PATH"), &ownerInfo.Ownerinfo{})

	beego.Router(beego.AppConfig.String("CREATE_CATEGORY_PATH"), &createCategory.CreateCategory{})
	beego.Router(beego.AppConfig.String("SEACRH_CATEGORY_PATH"), &searchCategory.SearchCategory{})
	beego.Router(beego.AppConfig.String("UPDATE_CATEGORY_PATH"), &updateCategory.UpdateCategory{})
	beego.Router(beego.AppConfig.String("VIEW_CATEGORY_PATH"), &viewCategory.ViewCategory{})

	beego.Router(beego.AppConfig.String("CREATE_SUBCATEGORY_PATH"), &createSubcategory.CreateSubcategory{})
	beego.Router(beego.AppConfig.String("SEACRH_SUBCATEGORY_PATH"), &searchSubcategory.SearchSubcategory{})
	beego.Router(beego.AppConfig.String("UPDATE_SUBCATEGORY_PATH"), &updateSubcategory.UpdateSubcategory{})
	beego.Router(beego.AppConfig.String("VIEW_SUBCATEGORY_PATH"), &viewSubcategory.ViewSubcategory{})

	beego.Router(beego.AppConfig.String("CUSTOMER_ORDER_BILL"), &customerorderBill.CustomerOrderBill{})
	beego.Router(beego.AppConfig.String("CLUSTER_PRODUCT_PURCHASE_BILL_PATH"), &clusterproductpurchaseBill.ClusterproductpurchaseBill{})

	beego.Router(beego.AppConfig.String("GET_GST"), &createProduct.GetGST{})

	beego.Router(beego.AppConfig.String("SEACRH_COMPLETE_DELIVERY_PATH"), &searchDelivery.SearchCompleteDelivery{})
	beego.Router(beego.AppConfig.String("UPDATE_COMPLETE_DELIVERY_PATH"), &updateDelivery.UpdateCompleteDelivery{})
	beego.Router(beego.AppConfig.String("VIEW_COMPLETE_DELIVERY_PATH"), &viewDelivery.ViewCompleteDelivery{})

	beego.Router(beego.AppConfig.String("SEACRH_REASSIGN_PATH"), &searchReassignOrder.SearchReassignOrder{})
	beego.Router(beego.AppConfig.String("UPDATE_REASSIGN__PATH"), &updateReassignOrder.UpdateReassignOrder{})
	beego.Router(beego.AppConfig.String("VIEW_REASSIGN__PATH"), &viewReassignOrder.ViewReassignOrder{})

	beego.Router(beego.AppConfig.String("CLUSTER_ADMIN_DASHBOARD"), &clusterdashbord.Clusterdashbord{})
	beego.Router(beego.AppConfig.String("CLUSTEREXECUTIVE_ADMIN_DASHBOARD"), &clusterexecutivedashbord.Clusterexecutivedashbord{})
	//beego.Router(beego.AppConfig.String("CLUSTEREXECUTIVE_ADMIN_DASHBOARD"), &hohelpdeskdashbord.Hohelpdeskdashbord{})

	beego.Router(beego.AppConfig.String("PRODUCER_BULK_UPLOAD_PATH"), &producerBulkUpload.ProducerBulkUpload{})
	beego.Router(beego.AppConfig.String("CUSTOMER_BULK_UPLOAD_PATH"), &customerBulkUpload.CustomerBulkUpload{})
	beego.Router(beego.AppConfig.String("CREATE_PRODUCERFORECAST_PATH"), &createProducerForecast.CreateProducerForecast{})
	beego.Router(beego.AppConfig.String("CREATE_CUSTOMERFORECAST_PATH"), &createCustomerForecast.CreateCustomerForecast{})
	beego.Router(beego.AppConfig.String("PRODUCERFORECAST_REPORT_PATH"), &producerForecastReport.ProducerForecastReport{})
	beego.Router(beego.AppConfig.String("CUSTOMERFORECAST_REPORT_PATH"), &customerForecastReport.CustomerForecastReport{})
	beego.Router(beego.AppConfig.String("FORECAST_DASHBOARD_PATH"), &forecastdashboard.ForecastDashboard{})
	beego.Router(beego.AppConfig.String("COMBINED_FORECAST_DASHBOARD_PATH"), &forecastdashboard.CombinedForecastDashboard{})

	beego.Router(beego.AppConfig.String("CREATE_SUPPLY_EXECUTIVE_PATH"), &createSupplyExecutive.CreateSupplyExecutive{})
	beego.Router(beego.AppConfig.String("SEACRH_SUPPLY_EXECUTIVE_PATH"), &searchSupplyExecutive.SearchSupplyExecutive{})
	beego.Router(beego.AppConfig.String("UPDATE_SUPPLY_EXECUTIVE_PATH"), &updateSupplyExecutive.UpdateSupplyExecutive{})
	beego.Router(beego.AppConfig.String("VIEW_SUPPLY_EXECUTIVE_PATH"), &viewSupplyExecutive.ViewSupplyExecutive{})

	beego.Router(beego.AppConfig.String("CREATE_CLUSTER_SUPPLYEXECUTIVE_PATH"), &createClusterSupplyexecutive.CreateClusterSupplyexecutive{})
	beego.Router(beego.AppConfig.String("SEACRH_CLUSTER_SUPPLYEXECUTIVE_PATH"), &searchClusterSupplyexecutive.SearchClusterSupplyexecutive{})
	beego.Router(beego.AppConfig.String("UPDATE_CLUSTER_SUPPLYEXECUTIVE_PATH"), &updateClusterSupplyexecutive.UpdateClusterSupplyexecutive{})
	beego.Router(beego.AppConfig.String("VIEW_CLUSTER_SUPPLYEXECUTIVE_PATH"), &viewClusterSupplyexecutive.ViewClusterSupplyexecutive{})

	beego.Router(beego.AppConfig.String("CREATE_CHANNEL_PARTNEER_TYPE_PATH"), &createChannelPartnerType.CreateChannelPartnerType{})
	beego.Router(beego.AppConfig.String("SEACRH_CHANNEL_PARTNEER_TYPE_PATH"), &searchChannelPartnerType.SearchChannelPartnerType{})
	beego.Router(beego.AppConfig.String("UPDATE_CHANNEL_PARTNEER_TYPE_PATH"), &updateChannelPartnerType.UpdateChannelPartnerType{})
	beego.Router(beego.AppConfig.String("VIEW_CHANNEL_PARTNEER_TYPE_PATH"), &viewChannelPartnerType.ViewChannelPartnerType{})

	beego.Router(beego.AppConfig.String("CREATE_CHANNEL_PARTNEER_PATH"), &createChannelPartner.CreateChannelPartner{})
	beego.Router(beego.AppConfig.String("SEACRH_CHANNEL_PARTNEER_PATH"), &searchChannelPartner.SearchChannelPartner{})
	beego.Router(beego.AppConfig.String("UPDATE_CHANNEL_PARTNEER_PATH"), &updateChannelPartner.UpdateChannelPartner{})
	beego.Router(beego.AppConfig.String("VIEW_CHANNEL_PARTNEER_PATH"), &viewChannelPartner.ViewChannelPartner{})

	beego.Router(beego.AppConfig.String("CREATE_CHANNEL_PARTNEER_CLUSTER_PATH"), &createChannelPartnerCluster.CreateChannelPartnerCluster{})
	beego.Router(beego.AppConfig.String("SEACRH_CHANNEL_PARTNEER_CLUSTER_PATH"), &searchChannelPartnerCluster.SearchChannelPartnerCluster{})
	beego.Router(beego.AppConfig.String("UPDATE_CHANNEL_PARTNEER_CLUSTER_PATH"), &updateChannelPartnerCluster.UpdateChannelPartnerCluster{})
	beego.Router(beego.AppConfig.String("VIEW_CHANNEL_PARTNEER_CLUSTER_PATH"), &viewChannelPartnerCluster.ViewChannelPartnerCluster{})

	beego.Router(beego.AppConfig.String("CREATE_CHANNEL_PARTNEER_PRODUCT_PATH"), &createChannelPartnerProduct.CreateChannelPartnerProduct{})
	beego.Router(beego.AppConfig.String("SEACRH_CHANNEL_PARTNEER_PRODUCT_PATH"), &searchChannelPartnerProduct.SearchChannelPartnerProduct{})
	beego.Router(beego.AppConfig.String("UPDATE_CHANNEL_PARTNEER_PRODUCT_PATH"), &updateChannelPartnerProduct.UpdateChannelPartnerProduct{})
	beego.Router(beego.AppConfig.String("VIEW_CHANNEL_PARTNEER_PRODUCT_PATH"), &viewChannelPartnerProduct.ViewChannelPartnerProduct{})

	beego.Router(beego.AppConfig.String("ORDER_TO_SUPPLYEXEC_PATH"), &orderToSupplyExecutive.OrderToSupplyExecutive{})
	beego.Router(beego.AppConfig.String("PICKUP_INVOICE_PATH"), &pickupinvoice.PickupInvoice{})

	beego.Router(beego.AppConfig.String("SURPLUS_MANAGEMENT_REPORT_PATH"), &surplusmanagementReport.SurplusManagementReport{})
	beego.Router(beego.AppConfig.String("SURPLUS_TRACKING_REPORT_PATH"), &surplusTrackingReport.SurplusTrackingReport{})

	beego.Router(beego.AppConfig.String("PAYU_RESPONSE_PATH"), &payuresponse.PayuResponse{})
	beego.Router(beego.AppConfig.String("PAYU_PAGE_PATH"), &payupage.PayUPage{})

}

