import React, { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../services/api';
import KanbanCard from '../kanbans/KanbanCard';

const CustomerDashboard = () => {
    const { customerId } = useParams();
    const navigate = useNavigate();
    const [kanbansByProduct, setKanbansByProduct] = useState({});
    const [availableCustomers, setAvailableCustomers] = useState([]);
    const [selectedCustomerId, setSelectedCustomerId] = useState(customerId || '');
    const [isLoading, setIsLoading] = useState(false);
    const [isDashboardDataReady, setIsDashboardDataReady] = useState(false); // ADD NEW isDashboardDataReady STATE


    useEffect(() => {
        fetchCustomersAndKanbans();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedCustomerId]);

    const fetchCustomersAndKanbans = async () => {
        setIsLoading(true);
        setIsDashboardDataReady(false); // **SET isDashboardDataReady to FALSE at start of fetch**
        try {
            const accountsResponse = await api.get('/accounts');
            // Filter for customers only
            const customers = accountsResponse.data;
            setAvailableCustomers(customers);

            const initialCustomerId = customerId || (customers.length > 0 ? customers[0].id : '');
            setSelectedCustomerId(initialCustomerId);


            if (initialCustomerId) {
                const kanbanResponse = await api.get(`/dashboards/customer/${selectedCustomerId}`);
                let kanbans = kanbanResponse.data.kanbans_by_product;
                kanbans = sortKanbansByCustomerSupplierAndDate(kanbans);
                setKanbansByProduct(kanbans);
            } else {
                setKanbansByProduct({});
            }


        } catch (error) {
            console.error("Error fetching data for customer dashboard", error);
        } finally {
            setIsLoading(false);
            setIsDashboardDataReady(true); // **SET isDashboardDataReady to TRUE in finally block** - Data loading is complete
        }
    };

    const handleCustomerChange = (e) => { // **CORRECT handleCustomerChange FOR CUSTOMER DASHBOARD**
        const newCustomerId = e.target.value;
        setSelectedCustomerId(newCustomerId);

        if (newCustomerId) {
            navigate(`/customer-dashboard/${newCustomerId}`); // Navigate to customer-dashboard/:customerId
        } else {
            navigate(`/customer-dashboard`); // Navigate to customer-dashboard
        }
    };

    const sortKanbansByCustomerSupplierAndDate = useCallback((kanbans) => {
        const sortedKanbans = { ...kanbans }; // Create a copy to avoid modifying original
        for (const product in sortedKanbans) {
            sortedKanbans[product].sort((a, b) => {
                const aMatch = a.customer_supplier === 2 ? -1 : 1; // Move matching to front
                const bMatch = b.customer_supplier === 2 ? -1 : 1; // Move matching to front
                if (aMatch !== bMatch) {
                    return aMatch - bMatch;
                }
                return new Date(a.data_aggiornamento) - new Date(b.data_aggiornamento); // Sort by date
            });
        }
        return sortedKanbans;
    }, []);

    const handleKanbanUpdate = useCallback((updatedKanban, productID) => {
        setKanbansByProduct(prevKanbansByProduct => {
            const updatedKanbansByProduct = { ...prevKanbansByProduct };
            if (updatedKanbansByProduct[productID]) {
                // **Simplified State Update - Replace Entire Kanban Object:**
                updatedKanbansByProduct[productID] = updatedKanbansByProduct[productID].map(k => {
                    if (k.kanban_id === updatedKanban.kanban_id) {
                        return updatedKanban; // **REPLACE ENTIRE KANBAN OBJECT with updatedKanban from API response**
                    } else {
                        return k;
                    }
                });
            }
            return updatedKanbansByProduct;
        });
    }, []);

    const handleStatusChangeSuccess = useCallback(() => {
        console.log("CustomerDashboard - handleStatusChangeSuccess CALLED - Re-fetching Kanban data");
        fetchCustomersAndKanbans(); // Re-fetch data from API
    }, [fetchCustomersAndKanbans]);

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

            {isLoading ? (
                <p>Loading Kanban data...</p>
            ) : (
                Object.keys(kanbansByProduct).length > 0 ? (
                    Object.keys(kanbansByProduct).map(product => (
                        <div key={product}>
                            <h3>{product}</h3>
                            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '20px' }}>
                            {kanbansByProduct[product].map(kanban => (
                            <KanbanCard
                                key={kanban.kanban_id}
                                kanban={kanban}
                                dashboardType="customer"
                                setKanbans={handleKanbanUpdate}
                                productID={product}
                                isDashboardDataReady={isDashboardDataReady}
                                onStatusChangeSuccess={handleStatusChangeSuccess} 
                            />
                        ))}
                            </div>
                        </div>
                    ))
                ) : (
                    <p>Please select a customer to view Kanban data</p>
                ))}


        </div>
    );
};

export default CustomerDashboard;