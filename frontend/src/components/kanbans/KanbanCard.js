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
                if (setKanbans) {
                    setKanbans(response.data, productID);
                }
            }
            if (onStatusChangeSuccess) {
                onStatusChangeSuccess();
            }


        } catch (error) {
            console.error("KanbanCard - Error updating Kanban status", error);
        }
    };

    return (
        <div className="kanban-card" style={{ borderColor: kanban.status_color, borderWidth: '2px' }}>
            <div className="kanban-card-header">
                <span className="kanban-card-id">#{kanban.kanban_id}</span>
                <h2>{kanban.status_name}</h2>
            </div>
            <div className="kanban-card-body">
                <p><strong>Product</strong> {kanban.product_name}</p>
                {dashboardType === 'customer' && kanban.supplier_name && (
                    <p><strong>Supplier</strong> {kanban.supplier_name}</p>
                )}
                {dashboardType === 'supplier' && kanban.customer_name && (
                    <p><strong>Customer</strong> {kanban.customer_name}</p>
                )}
                <p><strong>Container</strong> {kanban.tipo_contenitore}</p>
                <p><strong>Qty</strong> {kanban.quantity}</p>
                <p className="kanban-card-timestamp"><strong>Last Event</strong> {kanban.data_aggiornamento ? new Date(kanban.data_aggiornamento).toLocaleString() : '—'}</p>
            </div>
            {(dashboardType === 'supplier' && kanban.customer_supplier === 1) || (dashboardType === 'customer' && kanban.customer_supplier === 2) ? (
                <button
                    style={{ backgroundColor: kanban.status_color, color: 'white' }}
                    onClick={handleStatusChange}
                    disabled={!isDashboardDataReady}
                >
                    Change Status
                </button>
            ) : null}
        </div>
    );
};

export default KanbanCard;