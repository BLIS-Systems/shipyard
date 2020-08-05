import { Component, OnInit, Input } from '@angular/core';
import { UIControlGroup, ApiService, RoomConfig } from 'src/app/services/api.service';
import { MatDialog } from '@angular/material/dialog';
import { UIConfigDialogComponent } from './ui-config-dialog/ui-config-dialog.component';

@Component({
  selector: 'app-ui-config',
  templateUrl: './ui-config.component.html',
  styleUrls: ['./ui-config.component.scss']
})
export class UiConfigComponent implements OnInit {
  @Input('roomID') roomID: String;
  roomConf: RoomConfig;

  constructor(private api: ApiService,
    private dialog: MatDialog) {
    this.roomConf = this.api.getRoomConfig();
  }

  ngOnInit(): void {
  }

  addGroup(group: UIControlGroup, id: String) {
    var cgDialogRef;
    if (group == null) {
      cgDialogRef = this.dialog.open(UIConfigDialogComponent, {width: '60vw'});
    } else {
      cgDialogRef = this.dialog.open(UIConfigDialogComponent, {width: '60vw', data: {
        ControlGroup: group,
        ID: id
      }});
    }

    cgDialogRef.afterClosed().subscribe(data => {
      if (data != null) {
        this.roomConf.ControlGroups.set(data.ID, data.Config);
        if (id !== data.ID) {
          this.deleteGroup(id);
        }
      }
    });
  }

  deleteGroup(groupID: String) {
    this.roomConf.ControlGroups.delete(groupID);
  }

  getControlPanels(group: String): String[] {
    return ["one", "two", "three"];
  }

}
