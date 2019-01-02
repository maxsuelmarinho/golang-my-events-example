export interface Event {
    ID: string;
    Name: string;    
    StartDate: number;
    EndDate: number;    
    Location: {
        ID: string;
        Name: string;
        Address: string;
        Country: string;
        OpenTime: number;
        CloseTime: number;
    };    
}