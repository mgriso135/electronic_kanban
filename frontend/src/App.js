import React from 'react';
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import AccountList from './components/accounts/AccountList';
import AccountForm from './components/accounts/AccountForm';
import ProductList from './components/products/ProductList';
import ProductForm from './components/products/ProductForm';
import StatusList from './components/statuses/StatusList';
import StatusForm from './components/statuses/StatusForm';
import StatusChainList from './components/status_chains/StatusChainList';
import StatusChainForm from './components/status_chains/StatusChainForm';
import KanbanChainList from './components/kanban_chains/KanbanChainList';
import KanbanChainForm from './components/kanban_chains/KanbanChainForm';
import SupplierDashboard from './components/dashboards/SupplierDashboard';
import CustomerDashboard from './components/dashboards/CustomerDashboard';
import KanbanList from './components/kanbans/KanbanList'

function App() {
  return (
    <BrowserRouter>
        <nav>
            <Link to="/accounts">Accounts</Link> |
            <Link to="/products">Products</Link> |
          <Link to="/statuses">Statuses</Link> |
            <Link to="/status-chains">Status Chains</Link> |
            <Link to="/kanban-chains">Kanban Chains</Link> |
          <Link to="/kanbans">Kanbans</Link> |
            <Link to="/supplier-dashboard/1">Supplier Dashboard</Link>|
            <Link to="/customer-dashboard/1">Customer Dashboard</Link>
        </nav>
      <Routes>
          <Route path="/accounts" element={<AccountList />} />
          <Route path="/accounts/new" element={<AccountForm />} />
          <Route path="/accounts/:id/edit" element={<AccountForm />} />
          <Route path="/products" element={<ProductList />} />
          <Route path="/products/new" element={<ProductForm />} />
          <Route path="/products/:id/edit" element={<ProductForm />} />
          <Route path="/statuses" element={<StatusList />} />
          <Route path="/statuses/new" element={<StatusForm />} />
          <Route path="/statuses/:id/edit" element={<StatusForm />} />
          <Route path="/status-chains" element={<StatusChainList />} />
          <Route path="/status-chains/new" element={<StatusChainForm />} />
          <Route path="/status-chains/:id/edit" element={<StatusChainForm />} />
        <Route path="/kanban-chains" element={<KanbanChainList />} />
        <Route path="/kanban-chains/new" element={<KanbanChainForm />} />
        <Route path="/kanban-chains/:id/edit" element={<KanbanChainForm />} />
        <Route path="/kanbans" element={<KanbanList />} />
        <Route path="/supplier-dashboard/:supplierId" element={<SupplierDashboard />} />
        <Route path="/customer-dashboard/:customerId" element={<CustomerDashboard />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;