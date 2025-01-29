import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../services/api';
import KanbanCard from '../kanbans/KanbanCard';

const CustomerDashboard = () => {
    const { customerId } = useParams();
    const [kanbansByProduct, setKanbansByProduct] = useState({});
  const [availableCustomers, setAvailableCustomers] = useState([]);
  const [selectedCustomerId, setSelectedCustomerId] = useState(customerId || '');


    useEffect(() => {
    const fetchCustomersAndKanbans = async () => {
        try {
            const accountsResponse = await api.get('/accounts');
            // Filter for customers only
            const customers = accountsResponse.data;
            setAvailableCustomers(customers);


                if (selectedCustomerId) {
                    const kanbanResponse = await api.get(`/dashboards/customer/${selectedCustomerId}`);
                    let kanbans = kanbanResponse.data.kanbans_by_product;


                     for (const product in kanbans) {
                        kanbans[product].sort((a, b) => {
                             const aMatch = a.customer_supplier === 2 ? -1 : 1; // Move matching to front
                             const bMatch = b.customer_supplier === 2 ? -1 : 1; // Move matching to front
                             if (aMatch !== bMatch) {
                              return aMatch - bMatch;
                            }
                            return new Date(a.data_aggiornamento) - new Date(b.data_aggiornamento); // Sort by date
                        });
                     }
                    setKanbansByProduct(kanbans);

                }


        } catch (error) {
            console.error("Error fetching data for customer dashboard", error);
        }
    };
        fetchCustomersAndKanbans();
    }, [selectedCustomerId]);

    const handleCustomerChange = (e) => {
        setSelectedCustomerId(e.target.value);
    };

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

            {Object.keys(kanbansByProduct).length > 0 ? (
                Object.keys(kanbansByProduct).map(product => (
                    <div key={product}>
                        <h3>{product}</h3>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '20px'}}>
                         {kanbansByProduct[product].map(kanban => (
                            <KanbanCard key={kanban.kanban_id} kanban={kanban} dashboardType="customer" setKanbans={setKanbansByProduct} />
                         ))}
                    </div>
                  </div>
                ))
            ) : (
                <p>Please select a customer to view Kanban data</p>
            )}


        </div>
    );
};

export default CustomerDashboard;