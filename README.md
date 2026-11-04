# C2F Admin Web App

The Central Web Administration and ERP Portal for the Krishivala / C2F platform.

## Architecture & Features
- **Backend**: Go 1.23 + Beego MVC Framework
- **Frontend Views**: Server-side rendered HTML templates with AdminLTE / Bootstrap
- **Database**: AWS RDS MySQL (`C2F_NEW`)
- **Key Modules**:
  - Customer & Producer (Farmer) Management
  - Cluster Management & Local Delivery Charges
  - Product Catalog, Categories & GST Configuration
  - Order Routing, Reassignments & Delivery Confirmations
  - Collection Agent Cash Flow & Agent Assignments
  - Supply Chain Forecasting & Farmer Pickup Invoicing
  - RBAC (Role-Based Access Control) & System User Administration

## Running Locally
Run the convenient launcher script:
```cmd
start-admin.bat
```
Or with PowerShell:
```powershell
.\start-admin.ps1
```
The application will launch on **http://localhost:5005**.
