export interface Event {
    id: string;
    name: string;    
    startDate: number;
    endDate: number;    
    location: {
        id: string;
        name: string;
        address: string;
        country: string;
        openTime: number;
        closeTime: number;
    };    
}