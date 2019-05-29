import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { HttpHeaders} from '@angular/common/http'

const httpOptions = {
  headers: new HttpHeaders({
    'Content-Type':  'application/json',
    'Authorization': 'my-auth-token'
  })
};

@Injectable({
  providedIn: 'root'
})
export class HttpServiceService {

  url = 'http://mndkk.dk:50002';///boards/30aea474c680/ping';//Insert url to server

  constructor(private httpClient: HttpClient) { }

  public get(url:string): Observable<any>{
    return this.httpClient.get(this.url+url);
  }

  public post(url:string,data: any): Observable<any>{
    console.log(this.url + url +"POST" + data)
    return this.httpClient.post(this.url+url, data);
  }



}
 