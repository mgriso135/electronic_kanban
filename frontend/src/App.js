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
import KanbanForm from './components/kanbans/KanbanForm'

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
                <Link to="/supplier-dashboard">Supplier Dashboard</Link>| {/* Generic Supplier Dashboard Link */}
                <Link to="/customer-dashboard">Customer Dashboard</Link> {/* Generic Customer Dashboard Link */}
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
                <Route path="/kanbans" element={<KanbanList />} />
                <Route path="/kanbans/new" element={<KanbanForm />} /> {/* New route for creating Kanban */}
                <Route path="/kanbans/:id/edit" element={<KanbanForm />} /> {/* New route for editing Kanban */}
                <Route path="/supplier-dashboard" element={<SupplierDashboard />} /> {/* Generic path - no ID */}
                <Route path="/supplier-dashboard/:supplierId" element={<SupplierDashboard />} /> {/* Path with Supplier ID */}
                <Route path="/customer-dashboard" element={<CustomerDashboard />} /> {/* Generic path - no ID */}
                <Route path="/customer-dashboard/:customerId" element={<CustomerDashboard />} />   {/* Path with Customer ID */}
            </Routes>
        </BrowserRouter>
    );
}

export default App;