import React from 'react';
import api from '../../services/api';

const KanbanCard = ({ kanban, dashboardType, setKanbans, productID, isDashboardDataReady, onStatusChangeSuccess }) => { // ADD onStatusChangeSuccess PROP

    console.log("KanbanCard Component RENDERED - kanban_id:", kanban.kanban_id, "productID:", productID, "dashboardType:", dashboardType, "isDashboardDataReady:", isDashboardDataReady, "status_name:", kanban.status_name); // ADD THIS LOG - Log props on render

    const handleStatusChange = async () => {
        console.log("KanbanCard - handleStatusChange CALLED for kanban_id:", kanban.kanban_id, "productID:", productID);
        try {
            console.log("KanbanCard - handleStatusChange: Making API call...");
            const response = await api.put(`/kanbans/${kanban.kanban_id}/status`, { status_current: parseInt(kanban.status_current, 10) });
            console.log("KanbanCard - handleStatusChange: API response received:", response.data);

            if (dashboardType === "supplier" || dashboardType === "customer") {
                setKanbans(prevKanbansByProduct => {
                    console.log(`KanbanCard - setKanbans (${dashboardType} dashboard) CALLED with updatedKanban:`, response.data, "productID:", productID);

                    if (!prevKanbansByProduct[productID]) {
                        console.warn(`KanbanCard - setKanbans: productID "${productID}" not found in prevKanbansByProduct.`);
                        return prevKanbansByProduct || {};
                    }


                    const updatedKanbansByProduct = prevKanbansByProduct[productID].map(k => {
                        if (k.kanban_id === kanban.kanban_id) {
                            return { ...k, status_name: response.data.status_name, status_color: response.data.status_color, customer_supplier: response.data.customer_supplier, status_current: response.data.status_current };
                        } else {
                            return k;
                        }
                    });

                    // Immutably update the state - create a new object and assign the updated product kanbans
                    const nextKanbansByProduct = {
                        ...prevKanbansByProduct,
                        [productID]: updatedKanbansByProduct,
                    };

                    return nextKanbansByProduct; // Return the updated state
                }, () => {
                    console.log(`KanbanCard - setKanbans (${dashboardType} dashboard) CALLBACK EXECUTED`);
                    if (onStatusChangeSuccess) { // **CALL onStatusChangeSuccess CALLBACK HERE**
                        onStatusChangeSuccess(); // Force re-fetch in Dashboard
                    }
                });
            }


        } catch (error) {
            console.error("KanbanCard - Error updating Kanban status", error);
        }
    };

    return (
        <div className="kanban-card" style={{ borderColor: kanban.status_color, borderWidth: '3px' }}>
            <h2>{kanban.status_name}</h2>
            <p><strong>Product:</strong> {kanban.product_name}</p>
            {dashboardType === 'customer' && kanban.supplier_name && ( // Conditionally render Supplier Name in Customer Dashboard
                <p><strong>Supplier:</strong> {kanban.supplier_name}</p>
            )}
            {dashboardType === 'supplier' && kanban.customer_name && ( // Conditionally render Customer Name in Supplier Dashboard
                <p><strong>Customer:</strong> {kanban.customer_name}</p>
            )}
            <p><strong>Container:</strong> {kanban.tipo_contenitore}</p>
            <p><strong>Quantity:</strong> {kanban.quantity}</p>
            {(dashboardType === 'supplier' && kanban.customer_supplier === 1) || (dashboardType === 'customer' && kanban.customer_supplier === 2) ? ( // Conditional rendering for button
                <button
                    style={{ backgroundColor: kanban.status_color, color: 'white', padding: '10px' }}
                    onClick={handleStatusChange}
                    disabled={!isDashboardDataReady} // **DISABLE BUTTON INITIALLY, ENABLE CONDITIONALLY**
                >
                    Change Status
                </button>
            ) : null}
        </div>
    );
};

export default KanbanCard;