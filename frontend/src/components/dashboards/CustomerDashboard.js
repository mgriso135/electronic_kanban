import React, { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import api from '../../services/api';
import KanbanCard from '../kanbans/KanbanCard';

const CustomerDashboard = () => {
    const { customerId } = useParams();
    const navigate = useNavigate();
    const location = useLocation();
    const [kanbansByProduct, setKanbansByProduct] = useState({});
    const [availableCustomers, setAvailableCustomers] = useState([]);
    const [selectedCustomerId, setSelectedCustomerId] = useState(customerId || '');
    const [isLoading, setIsLoading] = useState(false); // Keep isLoading state INTERNALLY

    useEffect(() => {
        fetchCustomersAndKanbans();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedCustomerId]);

    const fetchCustomersAndKanbans = async () => {
        setIsLoading(true); // Set loading to true at start
        try {
            const accountsResponse = await api.get('/accounts');
            const customers = accountsResponse.data;
            setAvailableCustomers(customers);

            const initialCustomerId = customerId || (customers.length > 0 ? customers[0].id : '');
            setSelectedCustomerId(initialCustomerId);

            if (initialCustomerId) {
                const kanbanResponse = await api.get(`/dashboards/customer/${initialCustomerId}`);
                let kanbans = kanbanResponse.data.kanbans_by_product;
                kanbans = sortKanbansByCustomerSupplierAndDate(kanbans);
                setKanbansByProduct(kanbans);
            } else {
                setKanbansByProduct({});
            }
        } catch (error) {
            console.error("Error fetching data for customer dashboard", error);
        } finally {
            setIsLoading(false); // Set loading to false in finally block
        }
    };

    const handleCustomerChange = (e) => {
        const newCustomerId = e.target.value;
        setSelectedCustomerId(newCustomerId);

        if (newCustomerId) {
            navigate(`/customer-dashboard/${newCustomerId}`);
        } else {
            navigate(`/customer-dashboard`);
        }
    };

    const sortKanbansByCustomerSupplierAndDate = useCallback((kanbans) => {
        const sortedKanbans = { ...kanbans };
        for (const product in sortedKanbans) {
            sortedKanbans[product].sort((a, b) => {
                const aMatch = a.customer_supplier === 2 ? -1 : 1;
                const bMatch = b.customer_supplier === 2 ? -1 : 1;
                if (aMatch !== bMatch) {
                    return aMatch - bMatch;
                }
                return new Date(a.data_aggiornamento) - new Date(b.data_aggiornamento);
            });
        }
        return sortedKanbans;
    }, []);


    const handleKanbanUpdate = useCallback((updatedKanban, productID) => {
        // **ADD CONDITIONAL CHECK AT THE VERY BEGINNING:**
        if (!kanbansByProduct || !kanbansByProduct[productID] || !Array.isArray(kanbansByProduct[productID])) {
            console.warn("handleKanbanUpdate: Data not ready, aborting update. productID:", productID);
            return; // ABORT FUNCTION IF DATA IS NOT READY
        }


        setKanbansByProduct(prevKanbansByProduct => {
            // Safety checks (already present - keep them, but they might be redundant now)
            if (!prevKanbansByProduct || !prevKanbansByProduct[productID] || !Array.isArray(prevKanbansByProduct[productID])) {
                console.warn("handleKanbanUpdate: prevKanbansByProduct or product data is not properly initialized yet (inside setState). This should not happen frequently now.");
                return prevKanbansByProduct || {}; // Redundant safety return, but keep it.
            }

            const updatedProductKanbans = prevKanbansByProduct[productID].map(k =>
                k.kanban_id === updatedKanban.kanban_id ? updatedKanban : k
            );
            const nextKanbansByProduct = { ...prevKanbansByProduct };
            nextKanbansByProduct[productID] = sortKanbansByCustomerSupplierAndDate({ [productID]: updatedProductKanbans })[productID];
            return nextKanbansByProduct;
        });
    }, [sortKanbansByCustomerSupplierAndDate, kanbansByProduct]); // Keep kanbansByProduct in dependency array for now


    return (
        <div>
            <h2>Customer Dashboard</h2>
            <div>
                <label>Select Customer:</label>
                <select value={selectedCustomerId} onChange={handleCustomerChange} disabled={isLoading}>
                    <option value="">Select a Customer</option>
                    {availableCustomers.map(customer => (
                        <option key={customer.id} value={customer.id}>{customer.name}</option>
                    ))}
                </select>
            </div>

            {/* Conditional rendering for the entire Kanban card section */}
            {!isLoading && selectedCustomerId && Object.keys(kanbansByProduct).length > 0 ? (
                Object.keys(kanbansByProduct).map(product => (
                    <div key={product}>
                        <h3>{product}</h3>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '20px' }}>
                            {kanbansByProduct[product].map(kanban => (
                                <KanbanCard 
                                key={kanban.kanban_id} 
                                kanban={kanban} 
                                dashboardType="customer" 
                                setKanbans={handleKanbanUpdate}  // Function as setKanbans
                                productID={product}             // Product ID string as productID
                            />
                            ))}
                        </div>
                    </div>
                ))
            ) : (
                selectedCustomerId ? <p>No Kanban data found for this customer.</p> : <p>Please select a customer to view Kanban data</p>
            )}
        </div>
    );
};

export default CustomerDashboard;