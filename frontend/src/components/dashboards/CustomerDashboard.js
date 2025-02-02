import React, { useState, useEffect, useCallback } from 'react'; // Import useCallback
import { useParams, useNavigate, useLocation } from 'react-router-dom'; // Import useNavigate, useLocation
import api from '../../services/api';
import KanbanCard from '../kanbans/KanbanCard';

const CustomerDashboard = () => {
    const { customerId } = useParams();
    const navigate = useNavigate(); // Hook for navigation
    const location = useLocation(); // Hook to get current location
    const [kanbansByProduct, setKanbansByProduct] = useState({});
    const [availableCustomers, setAvailableCustomers] = useState([]);
    const [selectedCustomerId, setSelectedCustomerId] = useState(customerId || '');


    useEffect(() => {
        fetchCustomersAndKanbans();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedCustomerId]);

    const fetchCustomersAndKanbans = async () => {
        try {
            const accountsResponse = await api.get('/accounts');
            const customers = accountsResponse.data;
            setAvailableCustomers(customers);

            // If customerId is in URL, use it, otherwise default to first customer or empty selection
            const initialCustomerId = customerId || (customers.length > 0 ? customers[0].id : '');
            setSelectedCustomerId(initialCustomerId);


            if (initialCustomerId) { // Fetch kanbans only if a customer is selected
                const kanbanResponse = await api.get(`/dashboards/customer/${initialCustomerId}`);
                let kanbans = kanbanResponse.data.kanbans_by_product;
                kanbans = sortKanbansByCustomerSupplierAndDate(kanbans);
                setKanbansByProduct(kanbans);
            } else {
                setKanbansByProduct({}); // Clear kanbans if no customer selected
            }


        } catch (error) {
            console.error("Error fetching data for customer dashboard", error);
        }
    };

    const handleCustomerChange = (e) => {
        const newCustomerId = e.target.value;
        setSelectedCustomerId(newCustomerId);

        if (newCustomerId) {
            navigate(`/customer-dashboard/${newCustomerId}`); // Update URL to include customerId
        } else {
            navigate(`/customer-dashboard`); // Navigate to generic dashboard if no customer selected
        }

    };

    const sortKanbansByCustomerSupplierAndDate = useCallback((kanbans) => { // Use useCallback
        const sortedKanbans = { ...kanbans };
        for (const product in sortedKanbans) {
            sortedKanbans[product].sort((a, b) => {
                const aMatch = a.customer_supplier === 2 ? -1 : 1; // Customer role is 2
                const bMatch = b.customer_supplier === 2 ? -1 : 1; // Customer role is 2
                if (aMatch !== bMatch) {
                    return aMatch - bMatch;
                }
                return new Date(a.data_aggiornamento) - new Date(b.data_aggiornamento);
            });
        }
        return sortedKanbans;
    }, []);


    const handleKanbanUpdate = useCallback((updatedKanban, productID) => {
        setKanbansByProduct(prevKanbansByProduct => {
            console.log("handleKanbanUpdate CALLED for productID:", productID, "updatedKanban:", updatedKanban); // Log at start

            // Safety checks (keep these)
            if (!prevKanbansByProduct || !prevKanbansByProduct[productID] || !Array.isArray(prevKanbansByProduct[productID])) {
                console.warn("handleKanbanUpdate: prevKanbansByProduct or product data is not properly initialized yet.");
                return prevKanbansByProduct || {};
            }

            const updatedProductKanbans = prevKanbansByProduct[productID].map(k =>
                k.kanban_id === updatedKanban.kanban_id ? updatedKanban : k
            );

            // Create a COMPLETELY NEW kanbansByProduct object - deep copy to ensure change detection
            const nextKanbansByProduct = { ...prevKanbansByProduct };
            nextKanbansByProduct[productID] = sortKanbansByCustomerSupplierAndDate({ [productID]: updatedProductKanbans })[productID];

            // Create a brand new object for the entire kanbansByProduct state
            const trulyNewKanbansByProduct = {};
            for (const prodId in nextKanbansByProduct) {
                trulyNewKanbansByProduct[prodId] = [...nextKanbansByProduct[prodId]]; // Create new arrays for each product
            }


            console.log("handleKanbanUpdate - BEFORE setKanbansByProduct - nextKanbansByProduct:", trulyNewKanbansByProduct); // Log before setState

            return trulyNewKanbansByProduct; // Return the truly new object
        }, () => { // Add a setState callback function
            console.log("handleKanbanUpdate - setKanbansByProduct CALLBACK EXECUTED - state should be updated now."); // Log after setState is done
            console.log("Current kanbansByProduct state:", kanbansByProduct); // Log state *after* update (may be async, so might not be immediately updated here)
        });
    }, [sortKanbansByCustomerSupplierAndDate, kanbansByProduct]); // Include kanbansByProduct in dependency array - CAREFUL: potential infinite loop if not handled correctly. We use callback setState to mitigate.

    return (
        <div>
            <h2>Customer Dashboard</h2>
            <div>
                <label>Select Customer:</label>
                <select value={selectedCustomerId} onChange={handleCustomerChange}>
                    <option value="">Select a Customer</option>
                    {availableCustomers.map(customer => (
                        <option key={customer.id} value={customer.id}>{customer.name}</option>
                    ))}
                </select>
            </div>

            {selectedCustomerId && Object.keys(kanbansByProduct).length > 0 ? ( // Conditionally render kanban data
                Object.keys(kanbansByProduct).map(product => (
                    <div key={product}>
                        <h3>{product}</h3>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '20px' }}>
                            {kanbansByProduct[product].map(kanban => (
                                <KanbanCard
                                    key={kanban.kanban_id}
                                    kanban={kanban}
                                    dashboardType="customer"
                                    setKanbans={handleKanbanUpdate} // Pass handleKanbanUpdate
                                    productID={product}
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