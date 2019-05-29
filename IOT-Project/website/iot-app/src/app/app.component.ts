import { Component} from '@angular/core';
import { HttpServiceService} from 'src/app/http-service.service'


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers:[HttpServiceService]
})
export class AppComponent { 
  title = 'Pycom Information Dashboard';
  thumbLabel = true;
  value = 0; 
  max = 100;
  min = 0;
  checked = false;

  constructor(private httpService: HttpServiceService){}
  onChange(value){
    console.log(value)
    this.httpService.post('/boards/30aea474c680/ping','')//Send intensity value
      .subscribe()
  }

} 
